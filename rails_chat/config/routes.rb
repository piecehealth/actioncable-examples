Rails.application.routes.draw do
  # Define your application routes per the DSL in https://guides.rubyonrails.org/routing.html

  # Defines the root path route ("/")
  root "rooms#login"
  post "create_session", to: "rooms#create_session", as: "create_session"
  delete "destroy_session", to: "rooms#destroy_session", as: "destroy_session"

  get "/rooms", to: "rooms#index", as: "rooms"
  get "/rooms/:id", to: "rooms#show", as: "room"
end
