package server

import (
	"backend-crm/internal/clients"
	"backend-crm/internal/config"
	database "backend-crm/internal/database/mongodb"
	"backend-crm/internal/handler"
	"fmt"

	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type Server struct {
	cfg                                   config.HTTPServer
	customerHandler                       *handler.CustomerHandler
	vesselHandler                         *handler.VesselHandler
	supplierHandler                       *handler.SupplierHandler
	terminalHandler                       *handler.TerminalHandler
	agentHandler                          *handler.AgentHandler
	shipmentHandler                       *handler.ShipmentHandler
	feedHandler                           *handler.FeedHandler
	categoryManagementActivityTypeHandler *handler.CategoryManagementActivityTypeHandler
	categoryManagementProductTypeHandler  *handler.CategoryManagementProductTypeHandler
	checklistHandler                      *handler.ChecklistHandler

	loginHandler          *handler.LoginHandler
	invoicePricingHandler *handler.InvoicePricingHandler

	router *chi.Mux
}

func NewServer(
	cfg config.HTTPServer,
	customerCollection database.Collection[database.Customer, database.CustomerResponse],
	vesselCollection database.Collection[database.Vessel, database.VesselManagementResponse],
	supplierCollection database.Collection[database.Supplier, database.SupplierManagementResponse],
	terminalCollection database.Collection[database.Terminal, database.TerminalManagementResponse],
	agentCollection database.Collection[database.Agent, database.AgentManagementResponse],
	activityTypeCollection database.Collection[database.CategoryManagementActivityType, database.CategoryManagementActivityTypeResponse],
	productTypeCollection database.Collection[database.CategoryManagementProductType, database.CategoryManagementProductTypeResponse],
	shipmentCollection database.Collection[database.Shipment, database.ShipmentResponse],
	feedCollection database.Collection[database.FeedEmail, database.FeedEmailResponse],
	loginCollection database.Collection[database.TenantUser, database.TenantUserResponse],
	InvoicePricingCollection database.Collection[database.InvoicePricing, database.InvoicePricingResponse],
	sessionCollection database.Collection[database.Session, database.SessionResponse],
	checklistColection database.Collection[database.Checklist, database.ChecklistResponse],

) *Server {
	router := chi.NewRouter()

	// Configure CORS middleware to allow all origins
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"https://www.columbus-crm.com", "http://localhost:5173"}, // Allow all origins
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	router.Use(cors.Handler)

	customerHandler := handler.NewCustomerHandler(customerCollection)
	vesselHandler := handler.NewVesselHandler(vesselCollection)
	supplierHandler := handler.NewSupplierHandler(supplierCollection)
	terminalHandler := handler.NewTerminalHandler(terminalCollection)
	agentHandler := handler.NewAgentHandler(agentCollection)
	categoryManagementActivityTypeHandler := handler.NewActivityTypeHandler(activityTypeCollection)
	categoryManagementProductTypeHandler := handler.NewProductTypeHandler(productTypeCollection)
	shipmentHandler := handler.NewShipmentHandler(shipmentCollection)
	feedHandler := handler.NewFeedHandler(feedCollection, shipmentCollection, checklistColection)
	loginHandler := handler.NewLoginHandler(loginCollection, sessionCollection)
	invoicePricingHandler := handler.NewInvoicePricingHandler(InvoicePricingCollection)
	checklistHandler := handler.NewChecklistHandler(checklistColection)

	clients.InitClients()

	srv := &Server{
		cfg:                                   cfg,
		customerHandler:                       customerHandler,
		vesselHandler:                         vesselHandler,
		supplierHandler:                       supplierHandler,
		terminalHandler:                       terminalHandler,
		agentHandler:                          agentHandler,
		categoryManagementActivityTypeHandler: categoryManagementActivityTypeHandler,
		categoryManagementProductTypeHandler:  categoryManagementProductTypeHandler,
		shipmentHandler:                       shipmentHandler,
		feedHandler:                           feedHandler,
		loginHandler:                          loginHandler,
		invoicePricingHandler:                 invoicePricingHandler,
		checklistHandler:                      checklistHandler,
		router:                                router,
	}

	srv.routes()

	return srv
}

func (s *Server) Start(ctx context.Context) {
	server := http.Server{
		Addr:         fmt.Sprintf(":%d", s.cfg.Port),
		Handler:      s.router,
		IdleTimeout:  s.cfg.IdleTimeout,
		ReadTimeout:  s.cfg.ReadTimeout,
		WriteTimeout: s.cfg.WriteTimeout,
	}

	shutdownComplete := handleShutdown(func() {
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("server.Shutdown failed: %v\n", err)
		}
	})

	if err := server.ListenAndServe(); err == http.ErrServerClosed {
		<-shutdownComplete
	} else {
		log.Printf("http.ListenAndServe failed: %v\n", err)
	}

	log.Println("Shutdown gracefully")
}

func handleShutdown(onShutdownSignal func()) <-chan struct{} {
	shutdown := make(chan struct{})

	go func() {
		shutdownSignal := make(chan os.Signal, 1)
		signal.Notify(shutdownSignal, os.Interrupt, syscall.SIGTERM)

		<-shutdownSignal

		onShutdownSignal()
		close(shutdown)
	}()

	return shutdown
}
