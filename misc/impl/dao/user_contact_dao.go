/*
 * Copyright © 2019 Hedzr Yeh.
 */

package dao

import (
	"fmt"
	"github.com/hedzr/voxr-api/api/v10"
	"github.com/hedzr/voxr-api/models"
	"github.com/hedzr/voxr-common/dc"
	"github.com/hedzr/voxr-lite/internal/dbe"
	"github.com/hedzr/voxr-lite/internal/exception"
	"github.com/sirupsen/logrus"
	"strings"
)

func AddContactFromUser(uidOwner int64, uidFriend int64, groupName string) (res *v10.Contact, err error) {
	return AddContactFromUserInternal(uidOwner, uidFriend, groupName, 0, 0, "", "", "")
}

func AddContactFromUserInternal(uidOwner int64, uidFriend int64, groupName string, isGroup, relationship int32, remarkName, tags, remarks string) (res *v10.Contact, err error) {
	var (
		ok       bool
		user, uf *models.User
		cg       *models.ContactGroup
		rel      *models.ContactRelation
		c        *models.Contact
	)

	// tmpl = &models.User{Id:uidOwner,}
	user = &models.User{Id: uidOwner}
	uf = &models.User{Id: uidFriend}
	if ok, err = dbe.DBE.Engine().Get(user); ok && err != nil {
		err = exception.NewWith("AddContactFromUser 'find user' error", err)
		return
	}
	if ok, err = dbe.DBE.Engine().Get(uf); ok && err != nil {
		err = exception.NewWith("AddContactFromUser 'find friend' error", err)
		return
	}

	cg, err = ensureContactGroup(dbe.DBE, uidOwner, groupName)
	if err != nil || cg == nil {
		err = exception.NewWith("ensureContactGroup error", err)
		return
	}
	logrus.Debugf("    contact group ensured: %v", cg)

	c, err = saveOrGetContact(dbe.DBE, uf)
	if err != nil {
		err = exception.NewWith("saveOrGetContact error", err)
		return
	}

	rel, err = saveOrGetContactRelation(dbe.DBE, c, cg, int(relationship), remarkName, "", tags, remarks)
	if err != nil {
		err = exception.NewWith("saveOrGetContactRelation error", err)
		return
	}

	res = &v10.Contact{Cb: c.ToProto(), Relation: rel.ToProto(), Group: cg.ToProto(), User: uf.ToProto()}
	return
}

func GetContact(req *v10.GetContactReq) (res *v10.GetContactReply, err error) {
	res = &v10.GetContactReply{ProtoOp: v10.Op_GetContactAck, Seq: req.Seq}

	if req.ProtoOp != v10.Op_GetContact {
		err, res.ErrorCode = exception.New2(exception.UnknownRequest)
		return
	}

	var (
		rel  = &models.ContactRelation{UidOwner: req.UidOwner}
		c    = &models.Contact{Uid: req.UidFriend}
		cg   = &models.ContactGroup{UidOwner: req.UidOwner}
		rels []*models.ContactRelation
		ok   bool
	)

	err = dbe.DBE.Engine().Table(rel).Alias("r").
		Join("INNER", []string{c.TableName(), "c"}, "c.id=r.cid").
		Join("INNER", []string{cg.TableName(), "cg"}, "cg.id=r.gid and cg.uid_owner=r.uid_owner").
		Where("r.uid_owner=? and c.uid=?", req.UidOwner, req.UidFriend).Find(&rels)
	if err != nil {
		err, res.ErrorCode = exception.NewError2(exception.ContactRelationNotExist, err)
	}

	for _, r := range rels {
		cg.Id = r.Gid
		ok, err = dbe.DBE.Engine().Get(cg)
		if err != nil || !ok {
			err, res.ErrorCode = exception.NewError2(exception.ContactGroupNotExist, err)
		}

		c.Id = r.Cid
		ok, err = dbe.DBE.Engine().Get(c)
		if err != nil || !ok {
			err, res.ErrorCode = exception.NewError2(exception.ContactNotExist, err)
		}

		res = &v10.GetContactReply{ProtoOp: v10.Op_GetContactAck, Seq: req.Seq,
			UidOwner: req.UidOwner, UidFriend: req.UidFriend,
			Cb: c.ToProto(), Relation: r.ToProto(), Group: cg.ToProto()}
		return
	}
	return
}

func UpdateContact(req *v10.UpdateContactReq) (res *v10.UpdateContactReply, err error) {
	res = &v10.UpdateContactReply{ProtoOp: v10.Op_UpdateContactAck, Seq: req.Seq}

	if req.ProtoOp != v10.Op_UpdateContact {
		err = exception.New(exception.UnknownRequest)
		return
	}

	var (
		friend = &models.User{Id: req.UidFriend}
		c      = &models.Contact{Id: req.Cb.Id}
		cg     = &models.ContactGroup{Id: req.Group.Id, UidOwner: req.UidOwner, Name: req.Group.Name}
		rel    = &models.ContactRelation{Gid: req.Relation.Gid, Cid: req.Relation.Cid}
		ok     bool
		rows   int64
	)

	if c.Id == 0 || cg.Id == 0 || req.Cb.Id == 0 || req.Group.Id == 0 || req.Relation.Gid == 0 || req.Relation.Cid == 0 {
		cgid, cid, _, e := contactIdsBy(req.UidOwner, req.UidFriend)
		if e != nil {
			err = e
			return
		}
		req.Cb.Id = cid
		req.Cb.Uid = req.UidFriend
		c.Id = cid
		req.Group.Id = cgid
		cg.Id = cgid
		req.Relation.Gid = cgid
		req.Relation.Cid = cid
		req.Relation.UidOwner = req.UidOwner
		rel.Gid = cgid
		rel.Cid = cid
	}

	if ok, err = dbe.DBE.Engine().Get(friend); err != nil || !ok {
		err, res.ErrorCode = exception.NewError2(exception.UserNotExist, err)
		return
	}

	if ok, err = dbe.DBE.Engine().Get(c); err != nil || !ok {
		err, res.ErrorCode = exception.NewError2(exception.ContactNotExist, err)
		return
	}

	if ok, err = dbe.DBE.Engine().Get(cg); err != nil || !ok {
		err, res.ErrorCode = exception.NewError2(exception.ContactGroupNotExist, err)
		return
	}

	if ok, err = dbe.DBE.Engine().Get(rel); err != nil || !ok {
		err, res.ErrorCode = exception.NewError2(exception.ContactRelationNotExist, err)
		return
	}

	if err = dc.GormDefaultCopier.Copy(c, req.Cb); err != nil {
		err, res.ErrorCode = exception.NewError2(exception.UnknownError, err)
		return
	}

	if err = dc.GormDefaultCopier.Copy(cg, req.Group); err != nil {
		err, res.ErrorCode = exception.NewError2(exception.UnknownError, err)
		return
	}

	if err = dc.GormDefaultCopier.Copy(rel, req.Relation); err != nil {
		err, res.ErrorCode = exception.NewError2(exception.UnknownError, err)
		return
	}

	if len(c.FullName) == 0 {
		c.FullName = friend.FullName
	}
	if len(c.Nickname) == 0 {
		c.Nickname = friend.Nickname
	}
	if len(c.Tel) == 0 {
		c.Tel = friend.Mobile
	}
	if len(c.Email) == 0 {
		c.Email = friend.Email
	}
	if len(rel.RemarkName) == 0 {
		rel.RemarkName = friend.Nickname
	}
	if len(rel.RemarkMobile) == 0 {
		rel.RemarkMobile = friend.Mobile
	}
	if len(rel.RemarkEmail) == 0 {
		rel.RemarkEmail = friend.Email
	}

	if rows, err = dbe.DBE.Engine().Update(c, &models.Contact{Id: c.Id}); err != nil {
		err, res.ErrorCode = exception.NewError2(exception.CannotUpdateError, err)
		return
	} else {
		logrus.Debugf("contact: update %v row(s).", rows)
	}

	if rows, err = dbe.DBE.Engine().Update(cg, &models.ContactGroup{Id: req.Group.Id}); err != nil {
		err, res.ErrorCode = exception.NewError2(exception.CannotUpdateError, err)
		return
	} else {
		logrus.Debugf("contact group: update %v row(s).", rows)
	}

	if rows, err = dbe.DBE.Engine().Update(rel, &models.ContactRelation{Gid: req.Relation.Gid, Cid: req.Relation.Cid}); err != nil {
		err, res.ErrorCode = exception.NewError2(exception.CannotUpdateError, err)
		return
	} else {
		logrus.Debugf("contact relation: update %v row(s).", rows)
	}

	res.Ok = true
	return
}

func RemoveContact(req *v10.RemoveContactReq) (res *v10.RemoveContactReply, err error) {
	res = &v10.RemoveContactReply{ProtoOp: v10.Op_RemoveContactAck, Seq: req.Seq}
	if req.ProtoOp != v10.Op_RemoveContact {
		err, res.ErrorCode = exception.New2(exception.UnknownRequest)
		return
	}

	var (
		rows int64
		tmpl *models.ContactRelation
	)

	tmpl = &models.ContactRelation{
		Gid:      req.CgId,
		Cid:      req.CId,
		UidOwner: req.UidOwner,
	}

	rows, err = dbe.DBE.Engine().Delete(tmpl)
	if err != nil {
		err, res.ErrorCode = exception.NewError2("RemoveContact error", err)
		return
	}

	res.Ok = true
	logrus.Debugf("RemoveContact: %v rows deleted - relation[gid,cid,uidOwner]=[%v,%v,%v]", rows, req.CgId, req.CId, req.UidOwner)
	return
}

func ListContacts(uidOwner int64) (res *v10.ContactGroups, err error) {
	var (
		tmpl      *models.Contact
		relations []*models.ContactRelation
		groups    []*models.ContactGroup
		ok        bool
	)

	if err = dbe.DBE.Engine().Find(&groups, &models.ContactGroup{UidOwner: uidOwner}); err != nil {
		err, _ = exception.NewError2("ListContacts 'groups' error", err)
		return
	}

	res = &v10.ContactGroups{UidOwner: uidOwner, Groups: make([]*v10.ContactGroup, len(groups))}

	for ix, cg := range groups {
		res.Groups[ix] = &v10.ContactGroup{UidOwner: uidOwner, Cg: cg.ToProto()}

		relations = make([]*models.ContactRelation, 0)
		if err = dbe.DBE.Engine().Find(&relations, &models.ContactRelation{Gid: cg.Id, UidOwner: uidOwner}); err != nil {
			err, _ = exception.NewError2("ListContacts load relations error", err)
			return
		}

		var data = make([]*v10.ContactShort, 0)
		for _, r := range relations {
			tmpl = &models.Contact{Id: r.Cid}
			if ok, err = dbe.DBE.Engine().Get(tmpl); err != nil {
				err, _ = exception.NewError2("ListContacts load contact for relation error", err)
				return
			}
			if ok {
				data = append(data, &v10.ContactShort{Cb: tmpl.ToProto(), Relation: r.ToProto()})
			}
		}

		res.Groups[ix].Contacts = data
	}

	return
}

//

//

func saveOrGetContactRelation(db *dbe.DB, c *models.Contact, cg *models.ContactGroup, relationship int, remarkName, remarkAvatar, remarkTags, remakrs string) (rel *models.ContactRelation, err error) {
	var (
		ok   bool
		rows int64
		tmpl *models.ContactRelation
	)

	tmpl = &models.ContactRelation{
		Gid: cg.Id,
		// UidOwner: cg.UidOwner,
		Cid: c.Id,
	}
	ok, err = db.Engine().Unscoped().Exist(tmpl)
	if err != nil {
		return
	}

	if !ok {
		if len(remarkName) == 0 {
			remarkName = c.Nickname
		}

		rel = &models.ContactRelation{
			Gid:           cg.Id,
			UidOwner:      cg.UidOwner,
			Cid:           c.Id,
			RelationShip:  relationship,
			RemarkName:    remarkName,
			RemarkEmail:   c.Email,
			RemarkMobile:  c.Tel,
			RemarkTitle:   "", // TODO AddContactReq needs RemarkTitle
			RemarkOrgName: "", // TODO AddContactReq needs RemarkOrgName
			RemarkAvatar:  remarkAvatar,
			RemarkTags:    remarkTags,
			Remarks:       remakrs,
		}
		rows, err = db.Engine().Insert(rel)
		if err != nil {
			return
		}
		logrus.Debugf("    contact relation inserted: rows=%v, %v", rows, rel)
	} else {
		if r, e := db.Engine().Unscoped().Exec(fmt.Sprintf("UPDATE %v SET deleted_at=NULL where gid=? AND cid=?", rel.TableName()), cg.Id, c.Id); e != nil {
			logrus.Errorf("err: %v, r: %v", err, r)
			return
		}

		rel = &models.ContactRelation{
			RelationShip:  relationship,
			RemarkName:    remarkName,
			RemarkEmail:   c.Email,
			RemarkMobile:  c.Tel,
			RemarkTitle:   "", // TODO AddContactReq needs RemarkTitle
			RemarkOrgName: "", // TODO AddContactReq needs RemarkOrgName
			RemarkAvatar:  remarkAvatar,
			RemarkTags:    remarkTags,
			Remarks:       remakrs,
		}
		rows, err = db.Engine().Update(rel, tmpl)
		if err != nil {
			return
		}

		ok, err = db.Engine().Get(tmpl)
		if err != nil || !ok {
			return
		}
		rel = tmpl
		logrus.Debugf("    contact updated: rows=%v, %v", rows, rel)
	}
	return
}

func saveOrGetContact(db *dbe.DB, uf *models.User) (c *models.Contact, err error) {
	var (
		ok   bool
		rows int64
		tmpl *models.Contact
	)

	tmpl = &models.Contact{Uid: uf.Id}
	ok, err = db.Engine().Exist(tmpl)
	if err != nil {
		return
	}

	if !ok {
		c = &models.Contact{
			Uid:           uf.Id,
			Nickname:      uf.Nickname,
			FullName:      uf.FullName,
			Title:         "", // uf.Title
			OrgName:       "",
			Tel:           uf.Mobile,
			Email:         uf.Email,
			ImportRemarks: "",
		}
		rows, err = db.Engine().Insert(c)
		if err != nil {
			return
		}
		logrus.Debugf("    contact inserted: id=%v, rows=%v, %v", c.Id, rows, c)
	} else {
		ok, err = db.Engine().Get(tmpl)
		if err != nil || !ok {
			return
		}
		c = tmpl
		logrus.Debugf("    contact updated: id=%v, %v", c.Id, c)
	}
	return
}

func ensureDefaultContactGroup(db *dbe.DB, uidOwner int64) (res *models.ContactGroup, err error) {
	return ensureContactGroup(db, uidOwner, models.CGUnsorted)
}

func ensureContactGroup(db *dbe.DB, uidOwner int64, groupName string) (res *models.ContactGroup, err error) {
	var cc []*models.ContactGroup
	tmpl := &models.ContactGroup{UidOwner: uidOwner, Name: groupName}
	if err = db.Engine().Find(&cc, tmpl); err != nil {
		logrus.Errorf("Err: %v", err)
	} else {
		if len(cc) == 0 {
			logrus.Debugf("contact group not found, insert as new one.")
			res = tmpl // &models.ContactGroup{UidOwner: uidOwner, Name: groupName}
			var rows int64
			if rows, err = db.Engine().Insert(res); err != nil {
				logrus.Errorf("Err: %v", err)
			} else {
				logrus.Debugf("%d row(s) inserted: %v", rows, res)
			}
		} else {
			for ix, c := range cc {
				logrus.Debugf("contact group found: %3d, %v", ix, c)
				if strings.EqualFold(c.Name, groupName) {
					res = c
				}
			}
		}
	}
	return
}

func contactIdsBy(uidOwner, uidFriend int64) (cgid, cid int64, rels []*models.ContactRelation, err error) {
	var (
		rel = &models.ContactRelation{UidOwner: uidOwner}
		c   = &models.Contact{Uid: uidFriend}
		cg  = &models.ContactGroup{UidOwner: uidOwner}
	)

	err = dbe.DBE.Engine().Table(rel).Alias("r").
		Join("INNER", []string{c.TableName(), "c"}, "c.id=r.cid").
		Join("INNER", []string{cg.TableName(), "cg"}, "cg.id=r.gid and cg.uid_owner=r.uid_owner").
		Where("r.uid_owner=? and c.uid=?", uidOwner, uidFriend).Find(&rels)
	if err != nil {
		return
	}

	for _, r := range rels {
		cgid = r.Gid
		cid = r.Cid
		break
	}
	return
}
