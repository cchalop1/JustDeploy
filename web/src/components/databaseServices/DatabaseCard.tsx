import { Card } from "../ui/card";

export default function DatabaseCard() {
  return (
    <Card
      key={s.name}
      className="flex justify-between p-3 mb-3 hover:bg-gray-100"
    >
      <div className="flex flex-col gap-3">
        <div className="flex items-center gap-5">
          <img className="w-10" src={s.imageUrl}></img>
          <Status status={s.status} />
        </div>
      </div>
      <div>
        <Button variant="destructive" onClick={() => setServiceToDelete(s)}>
          Delete
        </Button>
      </div>
    </Card>
  );
}
