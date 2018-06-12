require 'socket'

s = TCPSocket.new 'localhost', 7890

module Redshift
  module OPC
    SystemId = 65535

    class Msg < Struct.new(:channel, :command, :length_high, :length_low, :data)
      def sysex?
        command == 255
      end

      def sysex_command
        if sysex?
          data[2]
        else
          nil
        end
      end

      def sysex_data
        if sysex?
          data.slice(4)
        else
          nil
        end
      end
    end

    class Conn
      def initialize(socket)
        @socket = socket
      end

      def send(bytes)
        @socket.write bytes.pack('C*')
      end

      def read
        header = @socket.read(4)
        data_length = header[2..3].map(&:chr).join.unpack('n')
        data = @socket.read(data_length) # make more robust
        # check length mismatch

        Msg.new(*header, data)
      end
    end
  end
end
