document.addEventListener('DOMContentLoaded', function() {
  // Add smooth scrolling for anchor links
  document.querySelectorAll('a[href^="#"]').forEach(anchor => {
    anchor.addEventListener('click', function(e) {
      e.preventDefault();
      
      const targetId = this.getAttribute('href');
      const targetElement = document.querySelector(targetId);
      
      if (targetElement) {
        window.scrollTo({
          top: targetElement.offsetTop - 80, // Offset for header
          behavior: 'smooth'
        });
      }
    });
  });

  // Add animation to cards on scroll
  const animateOnScroll = () => {
    const elements = document.querySelectorAll('.bento-item, .feature-card, .step-card');
    
    elements.forEach(element => {
      const elementTop = element.getBoundingClientRect().top;
      const elementVisible = 150;
      
      if (elementTop < window.innerHeight - elementVisible) {
        element.classList.add('animate');
      }
    });
  };

  // Add animation class to CSS
  const style = document.createElement('style');
  style.textContent = `
    .bento-item, .feature-card, .step-card {
      opacity: 0;
      transform: translateY(20px);
      transition: opacity 0.6s ease, transform 0.6s ease;
    }
    
    .animate {
      opacity: 1;
      transform: translateY(0);
    }
    
    .bento-item:nth-child(1) { transition-delay: 0.1s; }
    .bento-item:nth-child(2) { transition-delay: 0.2s; }
    .bento-item:nth-child(3) { transition-delay: 0.3s; }
    .bento-item:nth-child(4) { transition-delay: 0.4s; }
    
    .feature-card:nth-child(1) { transition-delay: 0.1s; }
    .feature-card:nth-child(2) { transition-delay: 0.2s; }
    .feature-card:nth-child(3) { transition-delay: 0.3s; }
    .feature-card:nth-child(4) { transition-delay: 0.4s; }
    .feature-card:nth-child(5) { transition-delay: 0.5s; }
    .feature-card:nth-child(6) { transition-delay: 0.6s; }
    
    .step-card:nth-child(1) { transition-delay: 0.1s; }
    .step-card:nth-child(2) { transition-delay: 0.2s; }
    .step-card:nth-child(3) { transition-delay: 0.3s; }
  `;
  document.head.appendChild(style);

  // Run animation on load and scroll
  window.addEventListener('scroll', animateOnScroll);
  animateOnScroll();

  // Add typing effect to code examples
  const codeBlocks = document.querySelectorAll('.code-container pre code');
  
  codeBlocks.forEach(codeBlock => {
    const originalContent = codeBlock.innerHTML;
    codeBlock.innerHTML = '';
    
    let i = 0;
    const typeCode = () => {
      if (i < originalContent.length) {
        codeBlock.innerHTML += originalContent.charAt(i);
        i++;
        setTimeout(typeCode, 5); // Fast typing speed
      }
    };
    
    // Start typing when code block is in view
    const observer = new IntersectionObserver(entries => {
      entries.forEach(entry => {
        if (entry.isIntersecting) {
          typeCode();
          observer.unobserve(entry.target);
        }
      });
    }, { threshold: 0.5 });
    
    observer.observe(codeBlock);
  });
});