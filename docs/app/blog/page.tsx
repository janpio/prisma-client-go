import { BlogClient } from "seobot";

import { ArticleCard } from "./card";
import { Pagination } from "./pagination";

import type { Metadata } from "next";

export const metadata: Metadata = {
	title: "Prisma Client Go Blog",
	description: "Blog posts about Prisma Client Go",
};

async function getPosts(page: number) {
	const key =
		process.env.SEOBOTAI_KEY || "a8c58738-7b98-4597-b20a-0bb1c2fe5772";

	const client = new BlogClient(key);
	return await client.getArticles(page, 10);
}

export const fetchCache = "force-no-store";

export default async function Blog({
	searchParams,
}: {
	searchParams: Promise<{ page: number }>;
}) {
	const { page } = await searchParams;
	const pageNumber = Math.max((page || 0) - 1, 0);
	const { total, articles } = await getPosts(pageNumber);
	const posts = articles || [];
	const lastPage = Math.ceil(total / 10);

	return (
		<section className="mx-auto my-8 max-w-3xl px-4 tracking-normal md:px-8 lg:mt-10 dark:text-white">
			<h1 className="my-4 text-4xl font-black">SeoBot Blog</h1>
			<ul>
				{posts.map((article: any) => (
					<ArticleCard key={article.id} article={article} />
				))}
			</ul>
			{lastPage > 1 && (
				<Pagination slug="/blog" pageNumber={pageNumber} lastPage={lastPage} />
			)}
		</section>
	);
}
