# Go-Restaurant-Management-Backend
Full backend for restaurant w/ Go + MongoDB + Gin + JWT Auth

This backend server handles routing, authentication w/ JWT tokens and database manipulation

## Installing
> [!IMPORTANT]    
> You must have [Go v1.21](https://go.dev/doc/install) (or higher) installed on your system
>
> You must have [MongoDB](https://www.mongodb.com/docs/v7.0/installation/) installed on your system

## Executing program

* Clone this repository to the location of your choosing
* Provide necessary env variables (i.e. *PORT* or *SECRET_KEY*) to *.env* file
* Provide necessary URI to *MongoDB* variable in *DBinstance* function located in *database/databaseConnection.go*
* Open your terminal
* Navigate to the saved location using ```cd folderName``` command, where *folderName* is the name of your path folder
* When in right location run:
```
go mod download
go build
go run main.go
```
## Usage (Requests)

* Viable operations with db (requests):
> User-related 
> ```
> /users - Get all user data from db (Method: GET)
> ```
> ```
> /users/:user_id - Get specified user by id data from db (Method: GET)
> ```
> ```
> /users/signup - Create new user w/ valid email, password and phone number (Method: POST)
> ```
> ```
> /users/login - Login with email and password (Method: POST)
> ```

> Menu-related
> ```
> /menus - Get all menu data from db (Method: GET)
> ```
> ```
> /menus/:menu_id - Get specified menu by id data from db (Method: GET)
> ```
> ```
> /menus - Create new menu w/ valid name and category (Method: POST)
> ```
> ```
> /menus/:menu_id - Update certain fields in specified menu (Method: PATCH)
> ```

> Food-related
> ```
> /foods - Get all food data from db (Method: GET)  // Provides pagination for the frontend
> ```
> ```
> /foods/:food_id - Get specified food by id data from db (Method: GET)
> ```
> ```
> /foods - Create new food item w/ valid name, price, image and menu_id (which menu this item belongs to)
> 
> (Method: POST)
> ```
> ```
> /foods/:food_id - Update certain fields in specified food item (Method: PATCH)
> ```

> Table-related
> ```
> /tables - Get all table data from db (Method: GET)
> ```
> ```
> /tables/:table_id - Get specified table by id data from db (Method: GET)
> ```
> ```
> /tables - Create new table entry w/ valid number of guests and table number (Method: POST)
> ```
> ```
> /tables/:table_id - Update certain fields in specified table entry (Method: PATCH)
> ```

> Order-related
> ```
> /orders - Get all order data from db (Method: GET)
> ```
> ```
> /orders/:order_id - Get specified order by id data from db (Method: GET)
> ```
> ```
> /orders - Create new order entry w/ valid order date and table_id (which table this order belongs to)
>
> (Method: POST)
> ```
> ```
> /orders/:order_id - Update certain fields in specified order entry (Method: PATCH)
> ```

> Ordered-items-related
> ```
> /orderItems - Get all ordered items data from db (Method: GET)
> ```
> ```
> /orderItems/:orderItem_id - Get specified ordered item by id data from db (Method: GET)
> ```
> ```
> /orderItems-order/:order_id - Get all ordered items data from db
>
> sorted by order_id (order items that belong to specified order) (Method: GET)
> ```
> ```
> /orderItems - Create new ordered items entry
>
> w/ valid quantity (small portion, medium or large), unit price, food_id (which food type these items belong to)
> and order_id (which order these items belong to)
>
> (Method: POST)
> ```
> ```
> /orderItems/:orderItem_id - Update certain fields in specified ordered items entry (Method: PATCH)
> ```

> Invoice-related
> ```
> /invoices - Get all invoice data from db (Method: GET)
> ```
> ```
> /invoices/:invoice_id - Get specified invoice by id data from db (Method: GET)
> ```
> ```
> /invoices - Create new invoice w/ valid payment method (card / cash / "") and payment status (pending / paid)
>
> (Method: POST)
> ```
> ```
> /invoices/:invoice_id - Update certain fields in specified invoice (Method: PATCH)
> ``` 

## Help

> [!NOTE]  
> For customizability of db items update the Structs provided in *models* folder accordingly (don't forget to update controller update functions as well)

> [!NOTE]  
> For extra field validation add the following line (change the field_name and validation_methods accordingly; validation_methods can be omitted) to the struct field after the data type:
> ```
> `json:"field_name" validate:"required,validation_methods"`
> ```

> [!IMPORTANT]  
> There must be no spaces between validation methods, required status and logical operators.


Feel free to report any issues or suggest improvements

## Version History

* v.0.0.1:

    * Initial Release

Tested w/ Postman
