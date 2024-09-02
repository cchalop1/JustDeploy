import { Card } from "../ui/card";
import { Skeleton } from "../ui/skeleton";

export default function DatabaseServiceCardLoading() {
  return (
    <Card className="flex justify-between p-3 h-24">
      <div className="flex flex-col gap-3 w-full">
        <div className="flex justify-between">
          <Skeleton className="w-44 h-6" />
          <Skeleton className="w-24 h-10" />
        </div>
        <Skeleton className="w-20 h-3" />
      </div>
    </Card>
  );
}
