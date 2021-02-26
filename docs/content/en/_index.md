---
title: K8ssandra
linkTitle: K8ssandra
---

<div id="home-header" class="container-fluid">
	<header class="row">
		<div class="col">
			<div class="container">
				<nav class="navbar navbar-expand-lg">
					<a class="navbar-brand" href="/"><img id="logo" src="/images/k8ssandra-stacked.svg" /><span class="sr-only">K8ssandra</span></a>
					<button class="navbar-toggler" type="button" data-toggle="collapse" data-target="#home-navigation-collapsible" aria-controls="home-navigation-collapsible" aria-expanded="false" aria-label="Toggle navigation">
						<i class="fas fa-bars"></i>
					</button>
					<div class="collapse navbar-collapse justify-content-end" id="home-navigation-collapsible">
						<ul class="navbar-nav">
							<li class="nav-item">
								<a class="nav-link" href="/about/">About</a>
							</li>
							<li class="nav-item">
								<a class="nav-link" href="/docs/">Documentation</a>
							</li>
							<li class="nav-item">
								<a class="nav-link" href="/community/">Community</a>
							</li>
						</ul>
						<form class="form-inline my-2 my-lg-0">
							<input type="search" class="form-control td-search-input" placeholder="&#xf002 Search this site…" aria-label="Search this site…" autocomplete="off">
						</form>
					</div>
				</nav>
				<div class="row">
					<div id="hero" class="col text-align-center">
						<div class="w-75 mx-auto text">
							K8ssandra provides a production-ready platform for running Apache Cassandra® on Kubernetes. This includes automation for operational tasks such as repairs, backups, and monitoring.
						</div>
						<div class="mx-auto">
							<a class="btn btn-lg btn-primary" href="{{< relref "docs" >}}">
								Learn More
							</a>
							<a class="btn btn-lg btn-secondary" href="https://github.com/k8ssandra/k8ssandra/releases">
								Download
							</a>
						</div>
					</div>
				</div>
			</div>
		</div>
	</header>
</div>

<div class="container">
	<main role="main" class="td-main">
		<div class="row">
			<section class="col">
				<div class="card text-center">
					<img src="/images/icons/helm.svg" />
					<h2>Helm</h2>
					<div class="description">
						Install the entire K8ssandra stack in <em>seconds</em> with Helm.<br /><br />
					</div>
					<div class="action">
						<a href="{{<relref "getting-started" >}}">Learn More</a>
					</div>
				</div>
			</section>
			<section class="col">
				<div class="card text-center">
					<img src="/images/icons/github.svg" />
					<h2>Contributions Welcome</h2>
					<div class="description">
						We follow the <a href="https://github.com/k8ssandra/k8ssandra/pulls">Pull Request</a> contributions workflow on <strong>GitHub</strong>. New users are always welcome!
					</div>
					<div class="action">
						<a href="https://github.com/k8ssandra/k8ssandra/pulls" target="_blank">Contribute</a>
					</div>
				</div>
			</section>
			<section class="col">
				<div class="card text-center">
					<img src="/images/icons/twitter.svg" />
					<h2>Follow Us on Twitter</h2>
					<div class="description">
						for announcements of latest features and releases.<br /><br />
					</div>
					<div class="action">
						<a href="https://twitter.com/k8ssandra">Follow @k8ssandra</a>
					</div>
				</div>
			</section>
		</div>
		<div class="row">
			<div class="col col-md-10 mx-auto quote">
				<blockquote>
					“New Relic is highly supportive of standardizing community-supported tools for operating and managing Cassandra clusters. We are excited about the K8ssandra launch and look forward to actively contributing and collaborating with the broader open source community. This is a great starting point for new and existing users to run Cassandra in Kubernetes and benefit from direct access to the best available Cassandra expertise and practices,”
				</blockquote>
				<cite>
				<strong>Tom Offermann</strong>, Lead Software Engineer at New Relic
				</cite>
			</div>
		</div>
		<div class="row">
			<div class="col text-center">
				<script id="asciicast-392352" src="https://asciinema.org/a/392352.js" async></script>
			</div>
		</div>
	</main>
</div>
