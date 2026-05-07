import {
  createFileRoute,
  isRedirect,
  redirect,
  useRouter,
} from "@tanstack/react-router";
import type { FormEvent } from "react";
import { Button } from "@/components/ui/Button";
import { Input } from "@/components/ui/Input";
import { getMe, useLogin } from "@/features/auth/api";

export const Route = createFileRoute("/auth/login")({
  beforeLoad: async () => {
    try {
      await getMe();
      throw redirect({ to: "/" });
    } catch (e) {
      if (isRedirect(e)) throw e;
    }
  },
  component: LoginPage,
});

function LoginPage() {
  const router = useRouter();
  const login = useLogin();

  function handleSubmit(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const data = new FormData(e.currentTarget);
    login.mutate(
      {
        data: {
          username: data.get("username") as string,
          password: data.get("password") as string,
        },
      },
      { onSuccess: () => router.navigate({ to: "/" }) },
    );
  }

  return (
    <form
      onSubmit={handleSubmit}
      className="mx-auto mt-16 flex max-w-sm flex-col gap-4"
    >
      <Input name="username" type="text" required placeholder="Username" />
      <Input name="password" type="password" required placeholder="Password" />
      <Button type="submit" disabled={login.isPending}>
        {login.isPending ? "Logging in…" : "Log in"}
      </Button>
      {login.isError && (
        <p className="text-sm text-red-500">Invalid credentials</p>
      )}
    </form>
  );
}
