"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { Button } from "@/components/ui/button";
import { useAuth } from "@/context/auth-context";
import { LogOut } from "lucide-react";

export default function Header() {
  const pathname = usePathname();
  const { user, logout } = useAuth();

  return (
    <header className="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
      <div className="container flex h-14 items-center">
        <Link href="/" className="flex items-center gap-2 font-semibold">
          <span>TaskMango</span>
        </Link>
        <nav className="ml-auto flex gap-4 sm:gap-6">
          {user ? (
            <>
              <Link
                href="/dashboard"
                className={`text-sm font-medium ${
                  pathname === "/dashboard"
                    ? "text-foreground"
                    : "text-muted-foreground"
                } transition-colors hover:text-foreground`}
              >
                Dashboard
              </Link>
              <Button
                variant="ghost"
                size="sm"
                className="gap-1"
                onClick={logout}
              >
                <LogOut className="h-4 w-4" />
                <span className="sr-only sm:not-sr-only sm:whitespace-nowrap">
                  Log Out
                </span>
              </Button>
            </>
          ) : (
            <>
              <Link
                href="/login"
                className={`text-sm font-medium ${
                  pathname === "/login"
                    ? "text-foreground"
                    : "text-muted-foreground"
                } transition-colors hover:text-foreground`}
              >
                Login
              </Link>
              <Link
                href="/register"
                className={`text-sm font-medium ${
                  pathname === "/register"
                    ? "text-foreground"
                    : "text-muted-foreground"
                } transition-colors hover:text-foreground`}
              >
                Register
              </Link>
            </>
          )}
        </nav>
      </div>
    </header>
  );
}
