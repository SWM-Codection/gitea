package ai

var (
	AiPullCommentService   PullCommentService
	AiPullCommentDbAdapter PullCommentDbAdapter
	AiPullCommentRequester PullCommentRequester
	AiSampleCodeService    SampleCodeService
	AiSampleCodeDbAdapter  SampleCodeDbAdapter
	AiSampleCodeRequester  SampleCodeRequester
)

func init() {
	AiPullCommentService = new(PullCommentServiceImpl)
	AiSampleCodeService = new(SampleCodeServiceImpl)
	AiSampleCodeDbAdapter = new(SampleCodeDbAdapterImpl)
	AiPullCommentDbAdapter = new(PullCommentDbAdapterImpl)
	AiPullCommentRequester = new(PullCommentRequesterImpl)
	AiSampleCodeRequester = new(SampleCodeRequesterImpl)
}
