import { cookies } from "next/headers";
import { NextResponse } from "next/server";

export async function POST() {
  // Clear the token cookie
  const cookieStore = await cookies();
  cookieStore.set({
    name: "token",
    value: "",
    httpOnly: true,
    path: "/",
    maxAge: 0, // Expire immediately
    sameSite: "strict",
  });

  return NextResponse.json({ success: true });
}
