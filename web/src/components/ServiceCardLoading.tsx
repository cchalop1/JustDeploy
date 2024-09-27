export function ServiceCardLoading() {
  return (
    <div className="w-80 h-36 bg-white border rounded-lg shadow-lg animate-pulse">
      <div className="flex justify-between ml-4 mr-4 items-center pt-3">
        <div className="w-8 h-8 bg-gray-200 rounded-full"></div>
        <div className="w-16 h-6 bg-gray-200 rounded"></div>
      </div>
      <div className="h-6 bg-gray-200 rounded mt-2 ml-4 mr-4 mb-2"></div>
      <div className="border-t h-2/6 flex items-center">
        <div className="ml-4 mr-4 mt-2 flex items-center">
          <div className="w-6 h-6 bg-gray-200 rounded-full"></div>
          <div className="w-24 h-4 bg-gray-200 rounded ml-2"></div>
        </div>
      </div>
    </div>
  );
}
