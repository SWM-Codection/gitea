package model

type ReactionList []*DiscussionReaction

func (list ReactionList) GroupByType() map[string]ReactionList {
	reactions := make(map[string]ReactionList)
	for _, reaction := range list {
		reactions[reaction.Type] = append(reactions[reaction.Type], reaction)
	}
	return reactions
}

func (list ReactionList) HasUser(userId int64) bool {
	if userId == 0 {
		return false
	}
	for _, reaction := range list {
		if reaction.UserId == userId {
			return true
		}
	}
	return false
}
