interface PaginationProps {
	slug: string;
	pageNumber: number;
	lastPage: number;
}

export const Pagination = ({ slug, pageNumber, lastPage }: PaginationProps) => {
	return (
		<div className="mt-12 flex items-center justify-center text-sm text-slate-300">
			<a
				className={`w-[90px] rounded-md border px-2 py-1 text-center ${pageNumber ? "" : "pointer-events-none opacity-30"}`}
				href={pageNumber ? `${slug}?page=${pageNumber}` : "#"}
			>
				← Prev
			</a>
			<div className="px-6 font-bold">
				{pageNumber + 1} / {lastPage}
			</div>
			<a
				className={`w-[90px] rounded-md border px-2 py-1 text-center ${
					pageNumber >= lastPage - 1 ? "pointer-events-none opacity-30" : ""
				}`}
				href={
					pageNumber >= lastPage - 1 ? "#" : `${slug}?page=${pageNumber + 2}`
				}
			>
				Next →
			</a>
		</div>
	);
};
