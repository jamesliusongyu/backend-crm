package server

import (
	// "backend-crm/pkg/auth"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func (s *Server) routes() {
	// Set content type middleware
	s.router.Use(render.SetContentType(render.ContentTypeJSON))

	// Public routes
	s.router.Route("/login", func(r chi.Router) {
		r.Post("/create_user", s.loginHandler.CreateTenantUser)
		r.Post("/", s.loginHandler.LoginAndAuth)
	})

	s.router.Route("/status_check", func(r chi.Router) {
		r.Get("/", s.loginHandler.StatusCheck)
	})

	// Private routes with JWT auth
	s.router.Group(func(r chi.Router) {
		// r.Use(auth.JWTMiddleware)
		r.Use(s.loginHandler.ValidateCookie)

		r.Route("/shipments", func(r chi.Router) {
			r.Get("/", s.shipmentHandler.GetAllShipment)
			r.Post("/", s.shipmentHandler.CreateShipment)
			r.Get("/{shipment_id}", s.shipmentHandler.GetShipmentFromId)
			r.Get("/query", s.shipmentHandler.FilterShipment)
			r.Delete("/{shipment_id}", s.shipmentHandler.DeleteShipmentById)
			r.Put("/{shipment_id}", s.shipmentHandler.UpdateShipmentById)
			r.Get("/statuses", s.shipmentHandler.GetAllShipmentStatuses)
			r.Get("/statuses_with_colours", s.shipmentHandler.GetAllShipmentStatusesWithColours)
			r.Get("/anchorage_locations", s.shipmentHandler.GetAllAnchorageLocations)
		})

		r.Route("/invoice", func(r chi.Router) {
			r.Get("/tenant", s.invoicePricingHandler.GetTenant)
			r.Get("/fees", s.invoicePricingHandler.GetInvoiceFeesFromPortAuthority)

			r.Get("/pda/{invoice_id}", s.invoicePricingHandler.GetPDAInvoicePricingFromId)
			r.Post("/pda", s.invoicePricingHandler.CreatePDAInvoicePricing)
			r.Put("/pda/{invoice_id}", s.invoicePricingHandler.EditPDAInvoicePricing)
		})

		r.Route("/feed", func(r chi.Router) {
			r.Get("/all", s.feedHandler.GetAllFeedEmails)
			r.Get("/", s.feedHandler.GetFeedEmailsByMasterEmail)
			r.Get("/{shipment_id}", s.feedHandler.GetFeedEmailsByShipmentId)
		})

		r.Route("/customer_management", func(r chi.Router) {
			r.Get("/", s.customerHandler.GetAllCustomers)
			r.Post("/", s.customerHandler.CreateCustomer)
			r.Get("/{customer_id}", s.customerHandler.GetCustomerFromId)
			r.Get("/{customer_name}", s.customerHandler.GetCustomerFromName)

			r.Get("/query", s.customerHandler.FilterCustomer)
			r.Delete("/{customer_id}", s.customerHandler.DeleteCustomerById)
			r.Put("/{customer_id}", s.customerHandler.UpdateCustomerById)
		})

		r.Route("/vessel_management", func(r chi.Router) {
			r.Get("/", s.vesselHandler.GetAllVessels)
			r.Post("/", s.vesselHandler.CreateVessel)
			r.Get("/{vessel_id}", s.vesselHandler.GetVesselById)
			r.Get("/query", s.vesselHandler.FilterVessel)
			r.Delete("/{vessel_id}", s.vesselHandler.DeleteVesselById)
			r.Put("/{vessel_id}", s.vesselHandler.UpdateVesselById)
		})

		r.Route("/supplier_management", func(r chi.Router) {
			r.Get("/", s.supplierHandler.GetAllSuppliers)
			r.Post("/", s.supplierHandler.CreateSupplier)
			r.Get("/{supplier_id}", s.supplierHandler.GetSupplierById)
			r.Get("/query", s.supplierHandler.FilterSupplier)
			r.Delete("/{supplier_id}", s.supplierHandler.DeleteSupplierById)
			r.Put("/{supplier_id}", s.supplierHandler.UpdateSupplierById)
		})

		r.Route("/terminal_management", func(r chi.Router) {
			r.Get("/", s.terminalHandler.GetAllTerminals)
			r.Post("/", s.terminalHandler.CreateTerminal)
			r.Get("/{terminal_id}", s.terminalHandler.GetTerminalById)
			r.Get("/query", s.terminalHandler.FilterTerminal)
			r.Delete("/{terminal_id}", s.terminalHandler.DeleteTerminalById)
			r.Put("/{terminal_id}", s.terminalHandler.UpdateTerminalById)
		})

		r.Route("/agent_management", func(r chi.Router) {
			r.Get("/", s.agentHandler.GetAllAgents)
			r.Post("/", s.agentHandler.CreateAgent)
			r.Get("/{agent_id}", s.agentHandler.GetAgentFromId)
			r.Get("/query", s.agentHandler.FilterAgent)
			r.Delete("/{agent_id}", s.agentHandler.DeleteAgentById)
			r.Put("/{agent_id}", s.agentHandler.UpdateAgentById)
		})

		r.Route("/category_management/activity_type", func(r chi.Router) {
			r.Get("/", s.categoryManagementActivityTypeHandler.GetAllCategoryManagementActivityTypes)
			r.Post("/", s.categoryManagementActivityTypeHandler.CreateActivityType)
			r.Get("/{activity_type_id}", s.categoryManagementActivityTypeHandler.GetActivityTypeFromId)
			r.Get("/activity_type_id", s.categoryManagementActivityTypeHandler.FilterActivityType)
			r.Delete("/{activity_type_id}", s.categoryManagementActivityTypeHandler.DeleteActivityTypeById)
		})

		r.Route("/category_management/product_type", func(r chi.Router) {
			r.Get("/", s.categoryManagementProductTypeHandler.GetAllCategoryManagementProductTypes)
			r.Get("/only_sub_products", s.categoryManagementProductTypeHandler.GetAllOnlySubProductTypes)
			r.Post("/", s.categoryManagementProductTypeHandler.CreateProductType)
			r.Get("/{product_type_id}", s.categoryManagementProductTypeHandler.GetProductTypeFromId)
			r.Get("/product_type_id", s.categoryManagementProductTypeHandler.FilterProductType)
			r.Delete("/{product_type_id}", s.categoryManagementProductTypeHandler.DeleteProductTypeById)
			r.Put("/{product_type_id}", s.categoryManagementProductTypeHandler.UpdateProductTypeById)
		})

		r.Route("/checklist", func(r chi.Router) {
			r.Get("/", s.checklistHandler.GetAllChecklist)
			r.Post("/", s.checklistHandler.CreateChecklist)
			r.Get("/{shipment_id}", s.checklistHandler.GetChecklistById)
			r.Get("/shipment_id", s.checklistHandler.FilterChecklist)
			r.Delete("/{shipment_id}", s.checklistHandler.DeleteChecklistById)
			r.Put("/{shipment_id}", s.checklistHandler.UpdateChecklistById)
		})

	})

	// Special routes without JWT auth
	s.router.Group(func(r chi.Router) {
		r.Route("/master_email_messages", func(r chi.Router) {
			r.Post("/", s.feedHandler.CreateFeedMessage)
		})
	})
}
