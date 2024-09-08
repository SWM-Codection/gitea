package dashboard

var (
	DpullDbAdapter    PullDbAdapter
	DdashboardService DashboardService
)

func init() {
	DpullDbAdapter = new(PullDbAdapterImpl)
	DdashboardService = new(DashboardServiceImpl)
}
