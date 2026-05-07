import { createFileRoute, useRouter } from "@tanstack/react-router";
import { Button } from "@/components/ui/Button";
import { useLogout } from "@/features/auth/api";

export const Route = createFileRoute("/")({
  component: HomePage,
});

function HomePage() {
  const router = useRouter();
  const logout = useLogout();

  function handleLogout() {
    logout.mutate(undefined, {
      onSuccess: () => router.navigate({ to: "/auth/login" }),
    });
  }

  return (
    <div className="mx-auto mt-16 flex max-w-sm flex-col gap-4">
      <p className="text-center">Home</p>
      <Button onClick={handleLogout} disabled={logout.isPending}>
        {logout.isPending ? "Logging out…" : "Log out"}
      </Button>
    </div>
  );
}
