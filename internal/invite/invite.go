package invite

import (
	"context"
	"log"

	"github.com/byteplow/idd4/internal/container"
	"github.com/google/uuid"

	keto "github.com/ory/keto-client-go"
)

var namespace_invite = "invites"
var relation_create_by = "create_by"
var relation_access = "access"
var relation_invited_by = "invited_by"

func CreateInvite(actorId string, endpoint string) (string, error) {
	invite := uuid.New().String()

	//set invites:invite#create@actorId
	_, _, err := container.KetoWriteClient.CreateRelationTuple(context.Background()).RelationQuery(keto.RelationQuery{
		Namespace: &namespace_invite,
		Object:    &invite,
		Relation:  &relation_create_by,
		SubjectId: &actorId,
	}).Execute()

	if err != nil {
		return "", err
	}

	//set invites:endpoint#access@invite
	_, _, err = container.KetoWriteClient.CreateRelationTuple(context.Background()).RelationQuery(keto.RelationQuery{
		Namespace: &namespace_invite,
		Object:    &endpoint,
		Relation:  &relation_access,
		SubjectId: &invite,
	}).Execute()

	if err != nil {
		return "", err
	}

	return invite, nil
}

func CheckInvite(invite string, endpoint string) bool {
	//check invites:endpoint#access@invite
	res, _, err := container.KetoReadClient.GetCheck(context.Background()).Namespace(namespace_invite).Object(endpoint).Relation(relation_access).SubjectId(invite).Execute()
	if err != nil {
		return false
	}

	return res.Allowed
}

func InvalidateInvite(invite string, actorId string) error {
	//get invites:invite#created@?
	res, _, err := container.KetoReadClient.GetRelationTuples(context.Background()).Namespace(namespace_invite).Object(invite).Relation(relation_create_by).Execute()
	if err != nil {
		return err
	}
	inviterId := res.RelationTuples[0].SubjectId

	//get invites:?#access@invite
	res, _, err = container.KetoReadClient.GetRelationTuples(context.Background()).Namespace(namespace_invite).SubjectId(invite).Relation(relation_access).Execute()
	if err != nil {
		return err
	}
	endpoint := res.RelationTuples[0].Object

	if actorId != "" {
		//set invites:inviter#invited_by@actorId
		_, _, err = container.KetoWriteClient.CreateRelationTuple(context.Background()).RelationQuery(keto.RelationQuery{
			Namespace: &namespace_invite,
			Object:    inviterId,
			Relation:  &relation_invited_by,
			SubjectId: &actorId,
		}).Execute()
		if err != nil {
			return err
		}

		//set invites:inviter#invited_by@(invites:actorId#invited_by)
		_, _, err = container.KetoWriteClient.CreateRelationTuple(context.Background()).RelationQuery(keto.RelationQuery{
			Namespace: &namespace_invite,
			Object:    inviterId,
			Relation:  &relation_invited_by,
			SubjectSet: &keto.SubjectSet{
				Namespace: namespace_invite,
				Object:    actorId,
				Relation:  relation_invited_by,
			},
		}).Execute()
		if err != nil {
			return err
		}
	}

	//del invites:invite#create@inviterId
	_, err = container.KetoWriteClient.DeleteRelationTuples(context.Background()).Namespace(namespace_invite).Object(invite).Relation(relation_create_by).SubjectId(*inviterId).Execute()
	if err != nil {
		return err
	}

	//del invites:endpoint#access@invite
	_, err = container.KetoWriteClient.DeleteRelationTuples(context.Background()).Namespace(namespace_invite).Object(endpoint).Relation(relation_access).SubjectId(invite).Execute()
	if err != nil {
		return err
	}

	return nil
}

func ListInvites(actorId string) ([]string, error) {
	invites := []string{}
	var nextPageToken *string

	for {
		//get invites:?#create@actorId
		req := container.KetoReadClient.GetRelationTuples(context.Background()).Namespace(namespace_invite).SubjectId(actorId).Relation(relation_create_by)

		if nextPageToken != nil {
			req.PageToken(*nextPageToken)
		}

		res, _, err := req.Execute()
		if err != nil {
			if err.Error() == "404 Not Found" {
				return invites, nil
			}
			return nil, err
		}

		for _, tuple := range res.RelationTuples {
			invites = append(invites, tuple.Object)
		}

		log.Println(*res.NextPageToken)
		if res.NextPageToken == nil || *res.NextPageToken == "" {

			return invites, nil
		}

		nextPageToken = res.NextPageToken
	}
}
