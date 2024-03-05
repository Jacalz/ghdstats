const std = @import("std");

pub fn main() !void {
    const args = std.os.argv;
    if (args.len == 0 or args.len > 2) {
        std.debug.print("Usage: gcdstats [user] [repository, optional]\n");
        return;
    }

    while (std.os.argv) |arg| {
        std.debug.print("{s}\n", .{arg});
    }
}
