package gitlab

import (
	"sort"
	"strings"

	log "github.com/sirupsen/logrus"
	gitlab "github.com/xanzy/go-gitlab"
)

func (c *Client) MailsFromGroupProject(group, project string) ([]string, error) {
	const level = 40
	// Works for gitlab 9, but documentation talks about https://docs.gitlab.com/ce/api/members.html#list-all-members-of-a-group-or-project-including-inherited-members
	// It doesn't work with curl + private token, and go-gitlab seems to not implement it
	groupMembers, resp, err := c.Groups.ListGroupMembers(group, &gitlab.ListGroupMembersOptions{})
	if err != nil {
		log.WithFields(
			log.Fields{
				"response": resp,
				"error":    err,
			},
		).Error("MailsFromGroupProject")
		return nil, err
	}
	mails := make(map[string]interface{})
	for _, member := range groupMembers {
		if member.AccessLevel < level {
			continue
		}
		user, resp, err := c.Users.GetUser(member.ID)
		if err != nil {
			log.WithFields(
				log.Fields{
					"response": resp,
					"error":    err,
				},
			).Error("MailsFromGroupProject")
			return nil, err
		}
		mails[user.Email] = true
	}

	id := strings.Join([]string{group, project}, "/")
	members, resp, err := c.ProjectMembers.ListProjectMembers(id, &gitlab.ListProjectMembersOptions{})
	if err != nil {
		log.WithFields(
			log.Fields{
				"response": resp,
				"error":    err,
			},
		).Error("MailsFromGroupProject")
		return nil, err
	}
	for _, member := range members {
		if member.State == "active" && member.AccessLevel >= level {
			mails[member.Email] = true
		}
	}
	smails := make([]string, len(mails))
	i := 0
	for key := range mails {
		smails[i] = key
		i++
	}
	sort.Strings(smails)
	return smails, nil
}
