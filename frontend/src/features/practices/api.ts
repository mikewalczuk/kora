import {
  useCreatePractice as _useCreatePractice,
  useGetPractice as _useGetPractice,
  useListPractices as _useListPractices,
} from "@/lib/api/generated/practices/practices";
import type {
  CreatePracticeResponse,
  ListPracticesParams,
  ListPracticesResponse,
  Practice,
} from "@/lib/api/generated/koraAPI.schemas";

export function useListPractices(params?: ListPracticesParams) {
  const { data, isLoading, isError } = _useListPractices(params);
  const response = data as unknown as ListPracticesResponse | undefined;
  return { practices: response?.items ?? [], total: response?.total ?? 0, isLoading, isError };
}

export function useCreatePractice() {
  const { mutateAsync, isPending } = _useCreatePractice();
  async function createPractice(noteId: string): Promise<CreatePracticeResponse> {
    const data = await mutateAsync({ data: { noteId } });
    return data as unknown as CreatePracticeResponse;
  }
  return { createPractice, isPending };
}

export function usePractice(id: string) {
  const { data, isLoading, isError } = _useGetPractice(id);
  const practice = data as unknown as Practice | undefined;
  return { practice, isLoading, isError };
}
