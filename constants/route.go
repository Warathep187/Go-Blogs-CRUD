package constants

type _RouteName struct {
	// auth
	LOGIN    string
	REGISTER string

	// blogs
	GET_BLOGS      string
	GET_BLOG_BY_ID string
	CREATE_BLOG    string
	UPDATE_BLOG    string
	DELETE_BLOG    string
}

var RouteName _RouteName

func init() {
	RouteName = _RouteName{
		// auth
		LOGIN:    "login",
		REGISTER: "register",

		// blogs
		GET_BLOGS:      "get_blogs",
		GET_BLOG_BY_ID: "get_blog_by_id",
		CREATE_BLOG:    "create_blog",
		UPDATE_BLOG:    "update_blog",
		DELETE_BLOG:    "delete_blog",
	}
}
