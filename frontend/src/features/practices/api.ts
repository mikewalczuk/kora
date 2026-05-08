import { useCreatePractice as _useCreatePractice, useGetPractice as _useGetPractice } from "@/lib/api/generated/practices/practices";
import type { Practice } from "@/lib/api/generated/koraAPI.schemas";

export function useCreatePractice() {
  const { mutateAsync, isPending } = _useCreatePractice();
  async function createPractice(noteId: string): Promise<Practice> {
    const data = await mutateAsync({ data: { noteId } });
    return data as unknown as Practice;
  }
  return { createPractice, isPending };
}

export function usePractice(id: string) {
  const { data, isLoading, isError } = _useGetPractice(id);
  const practice = data as unknown as Practice | undefined;
  return { practice, isLoading, isError };
}
