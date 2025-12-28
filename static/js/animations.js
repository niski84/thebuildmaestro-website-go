/**
 * Lightweight Animation System
 * Handles scroll reveals, parallax, and motion preferences
 */

(function() {
	'use strict';
	
	// Check for reduced motion preference
	const prefersReducedMotion = window.matchMedia('(prefers-reduced-motion: reduce)').matches;
	
	if (prefersReducedMotion) {
		// Disable all animations if user prefers reduced motion
		return;
	}
	
	// ============================================
	// Global effect 1: Scroll reveal + stagger
	// ============================================
	
	function initReveals() {
		const revealElements = document.querySelectorAll('.reveal');
		
		if (revealElements.length === 0) return;
		
		const observerOptions = {
			threshold: 0.1,
			rootMargin: '0px 0px -50px 0px'
		};
		
		const revealObserver = new IntersectionObserver((entries) => {
			entries.forEach(entry => {
				if (entry.isIntersecting) {
					entry.target.classList.add('is-visible');
					// Once visible, we can stop observing
					revealObserver.unobserve(entry.target);
				}
			});
		}, observerOptions);
		
		revealElements.forEach(el => {
			revealObserver.observe(el);
		});
	}
	
	// ============================================
	// Global effect 2: Parallax background
	// ============================================
	
	let parallaxSections = [];
	let rafId = null;
	let lastScrollY = 0;
	
	function initParallax() {
		parallaxSections = Array.from(document.querySelectorAll('.parallax-section')).map(section => {
			const bg = section.querySelector('.parallax-bg');
			if (!bg) return null;
			
			const speed = parseFloat(bg.dataset.speed) || 0.2;
			
			return {
				section,
				bg,
				speed,
				isVisible: false
			};
		}).filter(Boolean);
		
		if (parallaxSections.length === 0) return;
		
		// Use IntersectionObserver to track visibility
		const parallaxObserver = new IntersectionObserver((entries) => {
			entries.forEach(entry => {
				const sectionData = parallaxSections.find(s => s.section === entry.target);
				if (sectionData) {
					sectionData.isVisible = entry.isIntersecting;
				}
			});
		}, {
			threshold: 0,
			rootMargin: '50% 0px 50% 0px' // Start animating before fully visible
		});
		
		parallaxSections.forEach(({ section }) => {
			parallaxObserver.observe(section);
		});
		
		// Single scroll listener with requestAnimationFrame
		let ticking = false;
		
		function updateParallax() {
			const scrollY = window.pageYOffset || window.scrollY;
			const viewportHeight = window.innerHeight;
			
			parallaxSections.forEach(({ section, bg, speed, isVisible }) => {
				if (!isVisible) return;
				
				const rect = section.getBoundingClientRect();
				const sectionTop = rect.top + scrollY;
				const sectionHeight = rect.height;
				const scrolledPast = scrollY - sectionTop;
				
				// Calculate parallax offset - more accurate calculation matching original
				// The background moves at a different rate than the scroll
				// For layer-2, we want it to start below and move up into the gap
				const isLayer2 = bg.classList.contains('parallax-layer-2');
				let rate;
				
				if (isLayer2) {
					// Layer 2 starts at 100vh below, moves up much faster to create gap effect
					// The faster speed makes it catch up and appear in the gap between layers
					const initialOffset = viewportHeight;
					// Use a multiplier to make it move even faster relative to scroll
					rate = initialOffset - (scrolledPast * speed * 1.5);
				} else {
					// Layer 1 moves slower, creating depth
					rate = scrolledPast * speed;
				}
				
				bg.style.transform = `translateY(${rate}px)`;
			});
			
			ticking = false;
		}
		
		window.addEventListener('scroll', () => {
			if (!ticking) {
				window.requestAnimationFrame(updateParallax);
				ticking = true;
			}
		}, { passive: true });
		
		// Initial update
		updateParallax();
	}
	
	// ============================================
	// Initialize everything
	// ============================================
	
	function init() {
		initReveals();
		initParallax();
	}
	
	// Wait for DOM to be ready
	if (document.readyState === 'loading') {
		document.addEventListener('DOMContentLoaded', init);
	} else {
		init();
	}
	
})();

