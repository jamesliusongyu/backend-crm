package main

import (
	"backend-crm/internal/config"
	database "backend-crm/internal/database/mongodb"
	"backend-crm/internal/handler"
	"backend-crm/internal/server"
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// ctx := context.Background()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.Database.DatabaseURL))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	customerCollection := client.Database(cfg.Database.DatabaseName).Collection("CustomerManagementCollection")
	vesselCollection := client.Database(cfg.Database.DatabaseName).Collection("VesselManagementCollection")
	supplierCollection := client.Database(cfg.Database.DatabaseName).Collection("SupplierManagementCollection")
	terminalCollection := client.Database(cfg.Database.DatabaseName).Collection("TerminalManagementCollection")
	agentCollection := client.Database(cfg.Database.DatabaseName).Collection("AgentManagementCollection")
	categoryManagementActivityTypeCollection := client.Database(cfg.Database.DatabaseName).Collection("CategoryManagementActivityTypeCollection")
	categoryManagementProductTypeCollection := client.Database(cfg.Database.DatabaseName).Collection("CategoryManagementProductTypeCollection")
	shipmentCollection := client.Database(cfg.Database.DatabaseName).Collection("ShipmentsCollection")
	feedEmailCollection := client.Database(cfg.Database.DatabaseName).Collection("FeedEmailCollection")
	loginCollection := client.Database(cfg.Database.DatabaseName).Collection("LoginCollection")
	invoicePricingCollection := client.Database(cfg.Database.DatabaseName).Collection("InvoicePricingCollection")
	sessionCollection := client.Database(cfg.Database.DatabaseName).Collection("SessionCollection")
	checklistCollection := client.Database(cfg.Database.DatabaseName).Collection("ChecklistCollection")

	customerNewCollection := database.NewCustomerCollection(customerCollection)
	shipmentNewCollection := database.NewShipmentCollection(shipmentCollection)
	feedNewCollection := database.NewFeedMessageCollection(feedEmailCollection)
	vesselNewCollection := database.NewVesselCollection(vesselCollection)
	supplierNewCollection := database.NewSupplierCollection(supplierCollection)
	terminalNewCollection := database.NewTerminalCollection(terminalCollection)
	agentNewCollection := database.NewAgentCollection(agentCollection)
	categoryManagementActivityTypeNewCollection := database.NewActivityTypeCollection(categoryManagementActivityTypeCollection)
	categoryManagementProductTypeNewCollection := database.NewProductTypeCollection(categoryManagementProductTypeCollection)
	loginNewCollection := database.NewLoginCollection(loginCollection)
	invoicePricingNewCollection := database.NewInvoicePricingCollection(invoicePricingCollection)
	sessionNewCollection := database.NewSessionCollection(sessionCollection)
	checklistNewCollection := database.NewChecklistCollection(checklistCollection)

	server := server.NewServer(cfg.HTTPServer, customerNewCollection,
		vesselNewCollection, supplierNewCollection, terminalNewCollection, agentNewCollection, categoryManagementActivityTypeNewCollection, categoryManagementProductTypeNewCollection, shipmentNewCollection, feedNewCollection, loginNewCollection,
		invoicePricingNewCollection, sessionNewCollection, checklistNewCollection)

	stopChan := make(chan struct{})
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go handler.StartTimer(stopChan, wg, shipmentNewCollection, 2*time.Minute)

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Received shutdown signal")
		close(stopChan)
		cancel()
	}()
	server.Start(ctx)

	//The wg.Wait() call ensures that the application waits for all goroutines to finish before exiting.
	wg.Wait()
	log.Println("Application stopped")

}
