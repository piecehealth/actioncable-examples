class RoomChannel < ApplicationCable::Channel
  def subscribed
    if params[:id] == "forbidden_room" && current_user != "the chosen one"
      reject

      return
    end

    stream_from room_id
    Rails.logger.info("#{current_user} just join the #{room_id}.")
    ActionCable.server.broadcast(room_id, {user_join: "#{current_user} just join the room."})
  end

  def unsubscribed
    Rails.logger.info("#{current_user} just left the #{room_id}.")
    ActionCable.server.broadcast(room_id, {user_leave:  "#{current_user} just left the room."})
  end

  def send_message(data)
    ActionCable.server.broadcast(room_id, {send_by: current_user, message: data["message"]})
  end

  def kick(data)
    Rails.logger.info("disconnect #{data["name"]}'s connection.")
    ActionCable.server.remote_connections.where(current_user: data["name"]).disconnect
  end

  def stop_stream
    Rails.logger.info("stop all streams.")
    stop_all_streams
  end

  private
    def room_id
      "room_#{params[:id]}"
    end
end
