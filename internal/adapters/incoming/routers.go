package incoming

func RegisterCustomerRoutes(api huma.API, handler *Handler) {
	//huma.Register(api, huma.Operation{
	//	OperationID: "createproduct",
	//	Method:      "POST",
	//	Path:        "/Product",
	//	Summary:     "Create a Product",
	//}, handler.CreateCustomer)
	//
	//huma.Register(api, huma.Operation{
	//	OperationID: "listProduct",
	//	Method:      "GET",
	//	Path:        "/products",
	//	Summary:     "List all products",
	//}, handler.ListCustomers)
	//
	//huma.Register(api, huma.Operation{
	//	OperationID: "getproduct",
	//	Method:      "GET",
	//	Path:        "/product/{id}",
	//	Summary:     "Get a product by ID",
	//}, handler.GetCustomerByID)
	//
	//huma.Register(api, huma.Operation{
	//	OperationID: "getcustomerbyemail",
	//	Method:      "GET",
	//	Path:        "/customers/by-email",
	//	Summary:     "Get a customer by email",
	//}, handler.GetCustomerByEmail)
	//
	//huma.Register(api, huma.Operation{
	//	OperationID: "updatecustomer",
	//	Method:      "PUT",
	//	Path:        "/customers/{id}",
	//	Summary:     "Update a customer",
	//}, handler.UpdateCustomer)
	//
	//huma.Register(api, huma.Operation{
	//	OperationID: "deletecustomer",
	//	Method:      "DELETE",
	//	Path:        "/customers/{id}",
	//	Summary:     "Delete a customer",
	//}, handler.DeleteCustomer)

	huma.Register(api, huma.Operation{
		OperationID: "health",
		Method:      "GET",
		Path:        "/health",
		Summary:     "health is good",
	}, handler.ListCustomers)

}
