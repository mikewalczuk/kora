import { useQueryClient } from "@tanstack/react-query";
import { useListNotes as _useListNotes, useGetNote as _useGetNote, useCreateNote as _useCreateNote, useDeleteNote as _useDeleteNote, useUpdateNote as _useUpdateNote, getListNotesQueryKey, getGetNoteQueryKey } from "@/lib/api/generated/notes/notes";
import type { Note, ListNotesParams } from "@/lib/api/generated/koraAPI.schemas";

export function useNotesList(params?: ListNotesParams) {
  // customFetch returns the raw JSON body; the orval wrapper types don't apply at runtime
  const { data, isLoading, isError } = _useListNotes(params);
  const response = data as unknown as { items: Note[]; total: number } | undefined;
  return {
    notes: response?.items ?? [],
    total: response?.total ?? 0,
    isLoading,
    isError,
  };
}

export function useCreateNote() {
  const { mutateAsync, isPending } = _useCreateNote();
  async function createNote(title: string, content: string): Promise<Note> {
    const data = await mutateAsync({ data: { title, content } });
    return data as unknown as Note;
  }
  return { createNote, isPending };
}

export function useDeleteNote() {
  const queryClient = useQueryClient();
  const { mutateAsync, isPending } = _useDeleteNote();
  async function deleteNote(id: string) {
    await mutateAsync({ id });
    queryClient.invalidateQueries({ queryKey: getListNotesQueryKey() });
  }
  return { deleteNote, isPending };
}

export function useUpdateNote() {
  const queryClient = useQueryClient();
  const { mutateAsync, isPending } = _useUpdateNote();
  async function updateNote(id: string, title: string, content: string) {
    await mutateAsync({ id, data: { title, content } });
    queryClient.invalidateQueries({ queryKey: getGetNoteQueryKey(id) });
    queryClient.invalidateQueries({ queryKey: getListNotesQueryKey() });
  }
  return { updateNote, isPending };
}

export function useNote(id: string) {
  const { data, isLoading, isError } = _useGetNote(id);
  const note = data as unknown as Note | undefined;
  return { note, isLoading, isError };
}
