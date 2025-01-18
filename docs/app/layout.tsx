import Script from "next/script";
import { getPageMap } from "nextra/page-map";
import { Footer, Layout, Navbar } from "nextra-theme-docs";
import "nextra-theme-docs/style.css";

import "../globals.css";
import Analytics from "./analytics";

export const metadata = {
	title: {
		template: "%s – Prisma Client Go",
		default: "Prisma Client Go",
	},
	metadataBase: new URL("https://goprisma.org"),
	description:
		"Prisma Client Go is an auto-generated and fully type-safe database client",
	applicationName: "Go Prisma",
	icons: [
		{
			rel: "icon",
			url: "/icon.png",
			type: "image/png",
		},
	],
};

export default async function RootLayout({ children }) {
	return (
		<html lang="en" dir="ltr" suppressHydrationWarning>
			<head>
				<Script
					src="https://www.googletagmanager.com/gtag/js?id=G-SV698BGGW8"
					strategy="afterInteractive"
				/>
				<Script id="google-analytics" strategy="afterInteractive">
					{`
						window.dataLayer = window.dataLayer || [];
						function gtag(){dataLayer.push(arguments);}
						gtag('js', new Date());

						gtag('config', 'G-SV698BGGW8');
					`}
				</Script>
			</head>
			<body>
				<Analytics>
					<Layout
						footer={
							<Footer>
								<span>
									All source code and content licensed under&nbsp;
									<a
										href="https://github.com/steebchen/prisma-client-go/blob/main/LICENSE"
										target="_blank"
									>
										Apache 2.0
									</a>
								</span>
							</Footer>
						}
						navbar={
							<Navbar
								logo={<div>Go Prisma</div>}
								chatLink="https://discord.gg/er3ZbmYHDk"
								projectLink="https://github.com/steebchen/prisma-client-go"
							/>
						}
						editLink="Edit this page on GitHub"
						docsRepositoryBase="https://github.com/steebchen/prisma-client-go/tree/main/docs"
						sidebar={{ defaultMenuCollapseLevel: 1 }}
						pageMap={await getPageMap()}
					>
						{children}
					</Layout>
				</Analytics>
			</body>
		</html>
	);
}
