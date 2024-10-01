import SpinnerIcon from "@/assets/SpinnerIcon";
// import { Card, CardContent } from "../ui/card";

// export default function ServerListSkeleton() {
//   // Create an array of 4 items to represent loading skeleton cards
//   const skeletonCards = Array(4).fill(null);

//   return (
//     <div className="grid grid-cols-1 sm:grid-cols-1 md:grid-cols-1 lg:grid-cols-2 gap-3 h-2/3 mt-2">
//       {skeletonCards.map((_, index) => (
//         <Card className="h-full pt-4 pl-2 animate-pulse" key={index}>
//           <CardContent className="flex flex-col gap-2">
//             <div className="flex justify-between">
//               <div className="h-6 w-24 bg-gray-300 rounded"></div>
//               <div className="h-6 w-16 bg-gray-300 rounded"></div>
//             </div>
//             <div className="h-4 w-32 bg-gray-300 rounded"></div>
//             <div className="h-4 w-24 bg-gray-300 rounded"></div>
//           </CardContent>
//         </Card>
//       ))}
//     </div>
//   );
// }

export default function ServerListSkeleton() {
  return (
    <div className="flex justify-center items-center h-2/3">
      <SpinnerIcon color="text-black" />
    </div>
  );
}
