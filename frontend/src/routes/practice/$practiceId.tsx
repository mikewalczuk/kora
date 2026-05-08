import { createFileRoute, Link } from "@tanstack/react-router";
import { useMemo, useState } from "react";
import { useQueryClient } from "@tanstack/react-query";
import type { Exercise, MultiQuizQuestion, MultiQuizQuestionOptionsItem, Practice } from "@/lib/api/generated/koraAPI.schemas";
import { getGetPracticeQueryKey } from "@/lib/api/generated/practices/practices";
import { usePractice, useSubmitExercise } from "@/features/practices/api";

export const Route = createFileRoute("/practice/$practiceId")({
  component: PracticePage,
});

type QuestionResult = { selectedOptionId: string; correct: boolean };

function PracticePage() {
  const { practiceId } = Route.useParams();
  const { practice, isLoading, isError } = usePractice(practiceId);
  const queryClient = useQueryClient();

  function handleAllAnswered() {
    queryClient.invalidateQueries({ queryKey: getGetPracticeQueryKey(practiceId) });
  }

  return (
    <div className="mx-auto mt-16 max-w-lg px-4 flex flex-col gap-6">
      <div className="flex items-center justify-between">
        <Link to="/" className="text-sm text-gray-500 hover:text-gray-800">
          ← Back
        </Link>
        <span className="text-xs text-gray-400">Practice</span>
      </div>

      {isLoading && <p className="text-sm text-gray-400">Loading…</p>}
      {isError && <p className="text-sm text-red-500">Practice not found.</p>}
      {practice?.status === "completed" && <PracticeSummary practice={practice} />}
      {practice && practice.status !== "completed" && (
        <PracticeDetail practice={practice} onAllAnswered={handleAllAnswered} />
      )}
    </div>
  );
}

function PracticeSummary({ practice }: { practice: Practice }) {
  const questions = practice.exercises.flatMap((ex) => ex.questions);
  const total = questions.length;
  const correct = questions.filter((q) => q.submission?.correct).length;

  return (
    <div className="flex flex-col gap-4">
      <p className="text-sm font-semibold text-gray-800">Practice complete</p>
      <p className="text-2xl font-bold text-gray-900">
        {correct} <span className="text-gray-400 font-normal text-base">/ {total} correct</span>
      </p>
      <Link to="/" className="text-sm text-violet-600 hover:underline">
        Back to notes
      </Link>
    </div>
  );
}

function PracticeDetail({
  practice,
  onAllAnswered,
}: {
  practice: Practice;
  onAllAnswered: () => void;
}) {
  const totalQuestions = practice.exercises.reduce((sum, ex) => sum + ex.questions.length, 0);
  const [, setAnsweredCount] = useState(
    () => practice.exercises.reduce((sum, ex) => sum + ex.questions.filter((q) => q.submission != null).length, 0),
  );

  function handleQuestionAnswered() {
    setAnsweredCount((prev) => {
      const next = prev + 1;
      if (next >= totalQuestions) onAllAnswered();
      return next;
    });
  }

  return (
    <div className="flex flex-col gap-6">
      {practice.exercises.map((ex, i) => (
        <ExerciseCard
          key={ex.id}
          exercise={ex}
          practiceId={practice.id}
          index={i}
          onQuestionAnswered={handleQuestionAnswered}
        />
      ))}
    </div>
  );
}

function ExerciseCard({
  exercise,
  practiceId,
  index,
  onQuestionAnswered,
}: {
  exercise: Exercise;
  practiceId: string;
  index: number;
  onQuestionAnswered: () => void;
}) {
  const { submit, isPending } = useSubmitExercise();
  const initialResults = useMemo(() => {
    const map = new Map<string, QuestionResult>();
    for (const q of exercise.questions) {
      if (q.submission) {
        map.set(q.id, { selectedOptionId: q.submission.selectedOptionId, correct: q.submission.correct });
      }
    }
    return map;
  }, [exercise]);
  const [results, setResults] = useState<Map<string, QuestionResult>>(initialResults);

  async function handleSelect(questionId: string, selectedOptionId: string) {
    if (results.has(questionId) || isPending) return;
    const result = await submit(practiceId, exercise.id, questionId, selectedOptionId);
    setResults((prev) => new Map(prev).set(questionId, { selectedOptionId: result.selectedOptionId, correct: result.correct }));
    onQuestionAnswered();
  }

  return (
    <div className="flex flex-col gap-4">
      <p className="text-xs text-gray-400 uppercase tracking-wide">Exercise {index + 1}</p>
      {exercise.questions.map((q) => (
        <QuestionCard
          key={q.id}
          question={q}
          result={results.get(q.id) ?? null}
          onSelect={(optId) => handleSelect(q.id, optId)}
        />
      ))}
    </div>
  );
}

function QuestionCard({
  question,
  result,
  onSelect,
}: {
  question: MultiQuizQuestion;
  result: QuestionResult | null;
  onSelect: (optionId: string) => void;
}) {
  return (
    <div className="flex flex-col gap-3">
      <p className="text-sm font-medium text-gray-800">{question.question}</p>
      <ul className="flex flex-col gap-2">
        {question.options.map((opt) => (
          <OptionButton key={opt.id} option={opt} result={result} onSelect={onSelect} />
        ))}
      </ul>
      {result && (
        <p className={`text-xs font-medium ${result.correct ? "text-green-600" : "text-red-500"}`}>
          {result.correct ? "Correct!" : "Incorrect"}
        </p>
      )}
    </div>
  );
}

function OptionButton({
  option,
  result,
  onSelect,
}: {
  option: MultiQuizQuestionOptionsItem;
  result: QuestionResult | null;
  onSelect: (optionId: string) => void;
}) {
  const isSelected = result?.selectedOptionId === option.id;
  const answered = result !== null;

  let className = "w-full text-left rounded border px-3 py-2 text-sm transition-colors ";
  if (!answered) {
    className += "border-gray-200 text-gray-700 hover:border-violet-400 hover:bg-violet-50 cursor-pointer";
  } else if (isSelected && result.correct) {
    className += "border-green-400 bg-green-50 text-green-800 cursor-default";
  } else if (isSelected && !result.correct) {
    className += "border-red-400 bg-red-50 text-red-800 cursor-default";
  } else {
    className += "border-gray-100 text-gray-400 cursor-default";
  }

  return (
    <li>
      <button className={className} disabled={answered} onClick={() => onSelect(option.id)}>
        {option.text}
      </button>
    </li>
  );
}
