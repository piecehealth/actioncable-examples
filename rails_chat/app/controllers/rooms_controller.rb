class RoomsController < ApplicationController
  before_action :ensure_login, except: [:login, :create_session]

  helper_method :current_user

  def login
    redirect_to rooms_path if cookies.encrypted[:name].present?
  end

  def create_session
    cookies.encrypted[:name] = params[:name]

    redirect_to rooms_path
  end

  def destroy_session
    cookies.encrypted[:name] = nil

    redirect_to root_path
  end

  def index
  end

  def show
    @room_id = params[:id]
  end

  private
    def current_user
      @current_user ||= cookies.encrypted[:name]
    end

    def ensure_login
      redirect_to root_path if current_user.blank?
    end
end
