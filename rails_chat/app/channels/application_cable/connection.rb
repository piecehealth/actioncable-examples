module ApplicationCable
  class Connection < ActionCable::Connection::Base
    identified_by :current_user

    def connect
      self.current_user = find_verified_user
    end

    private
      def find_verified_user
        reject_unauthorized_connection if cookies.encrypted[:name].blank?

        cookies.encrypted[:name]
      end
  end
end
