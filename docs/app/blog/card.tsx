import Link from "next/link";

interface ArticleProps {
	article: any;
}

export const ArticleCard = ({ article }: ArticleProps) => {
	return (
		<li
			key={article.id}
			className="border-b border-gray-200 py-8 dark:border-slate-800"
		>
			<div className="flex w-full flex-wrap items-center gap-2 text-sm dark:text-slate-500">
				<span>
					Published{" "}
					{new Date(
						article.publishedAt || article.createdAt,
					).toLocaleDateString("en-US", {
						day: "numeric",
						month: "short",
						year: "numeric",
					})}
				</span>
				{article.readingTime ? (
					<span>{` ⦁ ${article.readingTime}`} min read</span>
				) : null}
			</div>
			<Link
				href={`/blog/${article.slug}`}
				className="mb-3 mt-2 block font-medium"
			>
				{article.headline}
			</Link>
			<div className="line mb-4 line-clamp-2 block text-sm text-slate-400 sm:text-base">
				{article.metaDescription}
			</div>
			<div className="flex flex-wrap justify-between gap-3">
				<div className="flex flex-wrap gap-2">
					{(article.tags || []).splice(0, 3).map((t: any, ix: number) => (
						<a
							key={ix}
							href={`/blog/tag/${t.slug}`}
							className="rounded bg-slate-800 px-2 py-1 text-xs text-slate-400"
						>
							{t.title}
						</a>
					))}
				</div>
				<Link
					href={`/blog/${article.slug}`}
					className="flex items-center text-sm font-medium text-orange-500 hover:text-orange-400"
				>
					Read More →
				</Link>
			</div>
		</li>
	);
};
