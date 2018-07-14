package actions

func (as *ActionSuite) Test_Admin_Index() {
	res := as.HTML("/admin").Get()
	as.Equal(200, res.Code)
	as.Contains(res.Body.String(), "this is the admin page")
}
