const std = @import("std");

pub fn main() !void {
    var args = std.process.args();

    while (args.next()) |arg| {
        std.debug.print("{s}\n", .{arg});
    }
}
