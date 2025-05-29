package repository

type Repositories interface {
	User() UserRepository
	UserSession() UserSessionRepository
	Post() PostRepository
	Comment() CommentRepository
	Category() CategoryRepository
	PostCategory() PostCategoryRepository
	PostReaction() PostReactionRepository
	CommentReaction() CommentReactionRepository
	PostAggregate() PostAggregateRepository
	UserAggregate() UserAggregateRepository
}
