require 'socket'

module Redshift
  module OPC
    SystemId = 65535

    module Cmd
      Welcome = 0

      OpenStream = 1
      CloseStream = 2
      SetStreamFps = 3

      SetEffectsJson = 4
      SetEffectsStreamFps = 5
      SetEffectsYaml = 6

      AppendEffectJson = 7
      AppendEffectYaml = 8

      OscSummary = 9
    end

    class Msg < Struct.new(:channel, :command, :data)
      def data
        super || self.data = []
      end

      def pixels
        if command == 0
          data.bytes.each_slice(3)
        end
      end

      def length_bytes
        [data.length].pack('n').bytes
      end

      def bytes
        [channel, command, *length_bytes, *data]
      end

      def self.sysex(channel, sysex_command, sysex_data=[])
        new(channel).tap { |msg|
          msg.sysex!
          msg.sysex_command = sysex_command
          msg.sysex_data = sysex_data
        }
      end

      def sysex?; command == 255; end
      def sysex_command; sysex? ? data[2] : nil; end
      def sysex_data; sysex? ? data[3..-1] : nil; end

      def sysex!
        self.command = 255

        system_id_high, system_id_low = [SystemId].pack('n').bytes
        self.data[0] = system_id_high
        self.data[1] = system_id_low
      end

      def sysex_command=(cmd)
        self.data[2] = cmd
      end

      def sysex_data=(data)
        self.data.insert(3, *data)
      end
    end

    class Conn
      def initialize(socket)
        @socket = socket
      end

      def read
        header = @socket.read(4)
        expected_data_length, _ = header[2..3].unpack('n')

        data = @socket.read(expected_data_length)
        raise Exception.new("data length mismatch") if data.length != expected_data_length

        Msg.new(header[0].ord, header[1].ord, data)
      end

      def send(msg)
        @socket.write msg.bytes.pack('C*')
      end

      def welcome
        send Msg.sysex(0, Cmd::Welcome)
        read
      end

      def open_stream(channel, desc)
        send Msg.sysex(channel, Cmd::OpenStream, desc.bytes)
      end

      def set_fps(channel, fps)
        send Msg.sysex(channel, Cmd::SetStreamFps, [fps])
      end

      def set_effects_json(channel, effects)
        send Msg.sysex(channel, Cmd::SetEffectsJson, effects.bytes)
      end

      def set_effects_yaml(channel, effects)
        send Msg.sysex(channel, Cmd::SetEffectsYaml, effects.bytes)
      end

      def append_effects_json(channel, effects)
        send Msg.sysex(channel, Cmd::AppendEffectsJson, effects.bytes)
      end

      def append_effects_yaml(channel, effects)
        send Msg.sysex(channel, Cmd::AppendEffectsYaml, effects.bytes)
      end

      def stream!(channel)
        loop {
          msg = read
          if msg.channel == channel and msg.command == 0
            yield msg.pixels
          end
        }
      end
    end
  end

  module Console
    class PixelPrinter < Struct.new(:output, :carriage_return, :line_feed)
      def print(pixels)
        output << "\r" if carriage_return
        pixels.each { |color| output << truecolor_string(color) }
        output << "\n" if line_feed
      end

      def truecolor_string(color, character=" ")
        r, g, b = color
        sprintf("\033[48;2;%d;%d;%dm%s\033[0m", r, g, b, character)
      end
    end

    def self.run
      socket = TCPSocket.new 'localhost', 7890
      conn = Redshift::OPC::Conn.new(socket)

      conn.welcome
      conn.open_stream(0, "strip")
      conn.set_fps(0, 60)

      printer = Redshift::Console::PixelPrinter.new(STDOUT, true, false)

      conn.stream!(0) { |pixels| printer.print(pixels) }
    end
  end
end

if __FILE__ == $0
  Redshift::Console.run
end
