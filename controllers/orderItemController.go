package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lackingworth/Go-Restaurant-Management/database"
	"github.com/lackingworth/Go-Restaurant-Management/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type OrderItemPack struct {
	Table_id		*string
	Order_items		[]models.OrderItem
}

var orderItemCollection *mongo.Collection = database.OpenCollection(database.Client, "orderItem")

func ItemsByOrder(id string) (OrderItems []primitive.M, err error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	// MongoDB Aggregation match stage
	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "order_id", Value: id}} /*end*/}}

	// MongoDB Aggregation stages for FoodID field
	lookupStage := bson.D{{Key: "$lookup", Value: bson.D{{Key: "from", Value: "food"}, {Key: "localField", Value: "food_id"}, {Key: "foreignField", Value: "food_id"}, {Key: "as", Value: "food"}} /*end*/}}
	unwindStage := bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$food"}, {Key: "preserveNullAndEmptyArrays", Value: true}} /*end*/}}

	// MongoDB Aggregation stages for OrderID field
	lookupOrderStage := bson.D{{Key: "$lookup", Value: bson.D{{Key: "from", Value: "order"}, {Key: "localField", Value: "order_id"}, {Key: "foreignField", Value: "order_id"}, {Key: "as", Value: "order"}} /*end*/}}
	unwindOrderStage := bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$order"}, {Key: "preserveNullAndEmptyArrays", Value: true}} /*end*/}}

	// MongoDB Aggregation stages for TableID field
	lookupTableStage := bson.D{{Key: "$lookup", Value: bson.D{{Key: "from", Value: "table"}, {Key: "localField", Value: "order.table_id"}, {Key: "foreignField", Value: "table_id"}, {Key: "as", Value: "table"}} /*end*/}}
	unwindTableStage := bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$table"}, {Key: "preserveNullAndEmptyArrays", Value: true}} /*end*/}}

	// MongoDB Aggregation 1st project stage
	projectStage := bson.D{{Key: "$project", Value: bson.D{{Key: "id", Value: 0}, {Key: "amount", Value: "$food.price"}, {Key: "total_count", Value: 1}, {Key: "food_name", Value: "$food.name"}, {Key: "food_image", Value: "$food.food_image"}, {Key: "table_number", Value: "$table.table_number"}, {Key: "table_id", Value: "$table.table_id"}, {Key: "order_id", Value: "$order.order_id"}, {Key: "price", Value: "$food.price"}, {Key: "quantity", Value: 1}} /*end*/}}

	// MongoDB Aggregation group stage
	groupStage := bson.D{{Key: "$group", Value: bson.D{{Key: "_id", Value: bson.D{{Key: "order_id", Value: "$order_id"}, {Key: "table_id", Value: "$table_id"}, {Key: "table_number", Value: "$table_number"}}}, {Key: "payment_due", Value: bson.D{{Key: "$sum", Value: "$amount"}}}, {Key: "total_count", Value: bson.D{{Key: "$sum", Value: 1}}}, {Key: "order_items", Value: bson.D{{Key: "$push", Value: "$$ROOT"}}} } /*end*/}}

	// MongoDB Aggregation 2nd project stage
	projectStage2 := bson.D{{Key: "$project", Value: bson.D{{Key: "id", Value: 0}, {Key: "payment_due", Value: 1}, {Key: "total_count", Value: 1}, {Key: "table_number", Value: "$_id.table_number"}, {Key: "order_items", Value: 1}} /*end*/}}

	res, err := orderCollection.Aggregate(ctx, mongo.Pipeline{
		matchStage,
		lookupStage,
		unwindStage,
		lookupOrderStage,
		unwindOrderStage,
		lookupTableStage,
		unwindTableStage,
		projectStage,
		groupStage,
		projectStage2,
	})

	if err != nil {
		panic(err)
	}

	if err := res.All(ctx, &OrderItems); err != nil {
		panic(err)
	}

	defer cancel()

	return OrderItems, err
}

func GetOrderItems() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		res, err := orderItemCollection.Find(context.TODO(), bson.M{})

		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error":"Error occured while listing ordered items"})
			return
		}

		var allOrderItems []bson.M

		if err := res.All(ctx, &allOrderItems); err != nil {
			log.Fatal(err)
			return
		}

		c.JSON(http.StatusOK, allOrderItems)
	}
}

func GetOrderItemsByOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		orderId := c.Param("order_id")

		allOrderItems, err := ItemsByOrder(orderId)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error":"Error occured while listing order items by order ID"})
			return
		}

		c.JSON(http.StatusOK, allOrderItems)
	}
}


func GetOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var orderItem models.OrderItem
		orderItemId := c.Param("order_item_id")

		err := orderItemCollection.FindOne(ctx, bson.M{"orderItem_id": orderItemId}).Decode(&orderItem)
		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error":"Error occured while listing single order item"})
			return
		}

		c.JSON(http.StatusOK, orderItem)
	}
}

func CreateOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var orderItemPack OrderItemPack
		var order models.Order
		defer cancel()

		if err := c.BindJSON(&orderItemPack); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		order.Order_Date, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		orderItemsToBeInserted := []interface{}{}
		order.Table_id = orderItemPack.Table_id
		order_id := OrderItemOrderCreator(order)

		for _, orderItem := range orderItemPack.Order_items {
			orderItem.Order_id = order_id

			validationErr := validate.Struct(orderItem)

			if validationErr != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
				return
			}

			orderItem.ID = primitive.NewObjectID()
			orderItem.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339)) 
			orderItem.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			orderItem.Order_item_id = orderItem.ID.Hex()

			var n = toFixed(*orderItem.Unit_price, 2)
			orderItem.Unit_price = &n
			orderItemsToBeInserted = append(orderItemsToBeInserted, orderItem)
		}

		res, insertErr := orderItemCollection.InsertMany(ctx, orderItemsToBeInserted)

		if insertErr != nil {
			log.Fatal(insertErr)
		}

		defer cancel()
		c.JSON(http.StatusOK, res)
	}
}

func UpdateOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var orderItem models.OrderItem
		var updateObj primitive.D
		orderItemId := c.Param("order_item_id")
		defer cancel()
		
		if orderItem.Unit_price != nil {
			updateObj = append(updateObj, bson.E{Key: "unit_price", Value: *&orderItem.Unit_price})
		}

		if orderItem.Quantity != nil {
			updateObj = append(updateObj, bson.E{Key: "quantity", Value: *orderItem.Quantity})
		}

		if orderItem.Food_id != nil {
			updateObj = append(updateObj, bson.E{Key: "food_id", Value: *orderItem.Food_id})
		}

		orderItem.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339)) 
		updateObj = append(updateObj, bson.E{Key: "updated_at", Value: orderItem.Updated_at})
		
		upsert := true
		filter := bson.M{"order_item_id": orderItemId}
		opt := options.UpdateOptions {
			Upsert: &upsert,
		}

		res, err := orderItemCollection.UpdateOne(
			ctx,
			filter,
			bson.D{{Key: "$set", Value: updateObj}},
			&opt,
		)

		if err != nil {
			msg := fmt.Sprintf("Order item update failed")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		defer cancel()
		c.JSON(http.StatusOK, res)
	}
}