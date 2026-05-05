import Button from "../../../shared/ui/Button";


export default function Pagination({ page, onPrev, onNext, disabledPrev, disabledNext }) {
  return (
    <div className="flex items-center justify-end gap-2 mt-4">
      <Button variant="secondary" onClick={onPrev} disabled={disabledPrev}>
        Prev
      </Button>
      <span className="text-sm text-gray-500">Page {page}</span>
      <Button variant="secondary" onClick={onNext} disabled={disabledNext}>
        Next
      </Button>
    </div>
  );
}