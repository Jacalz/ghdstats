const std = @import("std");

pub fn main() !void {
    const args = std.os.argv; // TODO: This does not support Windows!
    switch (args.len) {
        2, 3 => try display_statistics(args[1..]),
        else => {
            std.debug.print("Usage: gcdstats [user] [repository, optional]\n", .{});
            return;
        },
    }
}

const Repository = struct { name: []const u8 };

fn display_statistics(args: [][*:0]u8) !void {
    var repos: []Repository = undefined;

    if (args.len == 2) {
        repos = []Repository{ .name = args[0] + "/" + args[1] };
    } else {
        // TODO: Split on "/".
        repos = try fetch_repositories();
    }

    fetch_statistics(repos);
}

fn fetch_repositories() ![]Repository {}

fn fetch_statistics(_: []Repository) !void {}
