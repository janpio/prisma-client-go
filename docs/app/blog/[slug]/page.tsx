import Image from "next/image";
import Link from "next/link";
import { notFound } from "next/navigation";
import { BlogClient } from "seobot";

import "../blog.css";

import type { Metadata } from "next";

async function getPost(slug: string) {
	const key = process.env.SEOBOT_API_KEY;
	if (!key)
		throw Error(
			"SEOBOT_API_KEY enviroment variable must be set. You can use the DEMO key a8c58738-7b98-4597-b20a-0bb1c2fe5772 for testing - please set it in the root .env.local file",
		);

	const client = new BlogClient(key);
	return await client.getArticle(slug);
}

export const fetchCache = "force-no-store";

export async function generateMetadata({
	params: { slug },
}: {
	params: { slug: string };
}): Promise<Metadata> {
	const post = await getPost(slug);
	if (!post) return {};

	const title = post.headline;
	const description = post.metaDescription;
	return {
		title,
		description,
		metadataBase: new URL("https://devhunt.org"),
		alternates: {
			canonical: `/blog/${slug}`,
		},
		openGraph: {
			type: "article",
			title,
			description,
			images: [post.image],
			url: `https://devhunt.org/blog/${slug}`,
		},
		twitter: {
			title,
			description,
			card: "summary_large_image",
			images: [post.image],
		},
	};
}

export default async function Article({
	params: { slug },
}: {
	params: { slug: string };
}) {
	const post = await getPost(slug);
	if (!post) {
		return notFound();
	}

	console.log("post", post); // todo remove

	return (
		<section className="mx-auto max-w-2xl px-4 md:px-8 lg:my-8 dark:text-white">
			{post.category ? (
				<div className="mb-1 flex w-full flex-wrap items-center gap-2 text-sm dark:text-slate-400">
					<a href="/">Home</a>
					<svg
						width="12"
						height="12"
						viewBox="0 0 1024 1024"
						xmlns="http://www.w3.org/2000/svg"
					>
						<path
							fill="currentColor"
							d="M338.752 104.704a64 64 0 0 0 0 90.496l316.8 316.8l-316.8 316.8a64 64 0 0 0 90.496 90.496l362.048-362.048a64 64 0 0 0 0-90.496L429.248 104.704a64 64 0 0 0-90.496 0z"
						/>
					</svg>
					<Link href="/blog/">Blog</Link>
					<svg
						width="12"
						height="12"
						viewBox="0 0 1024 1024"
						xmlns="http://www.w3.org/2000/svg"
					>
						<path
							fill="currentColor"
							d="M338.752 104.704a64 64 0 0 0 0 90.496l316.8 316.8l-316.8 316.8a64 64 0 0 0 90.496 90.496l362.048-362.048a64 64 0 0 0 0-90.496L429.248 104.704a64 64 0 0 0-90.496 0z"
						/>
					</svg>
					<Link href={`/blog/category/${post.category.slug}`}>
						{post.category.title}
					</Link>
				</div>
			) : null}
			<div className="flex w-full flex-wrap items-center gap-2 text-sm dark:text-slate-400">
				<span>
					Published{" "}
					{new Date(post.publishedAt || post.createdAt).toLocaleDateString(
						"en-US",
						{
							day: "numeric",
							month: "short",
							year: "numeric",
						},
					)}
				</span>
				{post.readingTime ? (
					<span>{` ⦁ ${post.readingTime}`} min read</span>
				) : null}
			</div>
			<div className="relative mt-2 flex aspect-video w-full items-center justify-center overflow-hidden rounded-xl text-center">
				<Image
					src={post.image}
					alt={post.headline}
					layout="fill"
					objectFit="cover"
					className="!inset-auto"
				/>
			</div>
			<div
				className="prose prose-h1:text-slate-100 prose-h2:text-slate-100 prose-h3:text-slate-100 prose-strong:text-slate-100 mt-8 dark:text-slate-100"
				dangerouslySetInnerHTML={{ __html: post.html }}
			/>
			<div className="flex w-full flex-wrap justify-start gap-2">
				{(post.tags || []).map((t: any, ix: number) => (
					<a
						key={ix}
						href={`/blog/tag/${t.slug}`}
						className="rounded px-3 text-sm dark:bg-slate-700"
					>
						{t.title}
					</a>
				))}
			</div>
			{post.relatedPosts?.length ? (
				<div>
					<h2>Related posts</h2>
					<ul className="text-base">
						{post.relatedPosts.map((p: any, ix: number) => (
							<li key={ix}>
								<a href={`/blog/${p.slug}`}>{p.headline}</a>
							</li>
						))}
					</ul>
				</div>
			) : null}
		</section>
	);
}
