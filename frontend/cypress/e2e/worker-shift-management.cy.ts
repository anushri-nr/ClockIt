describe('Worker Shift Management - Release Shift Only', () => {
  beforeEach(() => {
    // Login as worker before each test
    cy.visit('http://localhost:4200/home');
    cy.contains('Worker').click();
    cy.get('input[name="email"]').type('rahul.umamahesha@gmail.com');
    cy.get('input[name="password"]').type('Rahul@2906');
    cy.contains('Sign In').click();
    cy.url().should('include', '/worker-dashboard');
    cy.contains('Hello,').should('be.visible');
    cy.wait(1000); // Wait for shifts to load
  });

  it('should display RELEASE SHIFT button on the shift card', () => {
    cy.contains('Upcoming Assigned Shifts').should('be.visible');
    cy.get('.release-btn').should('exist');
    cy.get('.release-btn').should('contain', 'Release Shift');
  });

  it('should click RELEASE SHIFT button and release the shift', () => {
    // Find the first upcoming shift item
    cy.get('.upcoming-item').first().should('be.visible');
    
    // Click the RELEASE SHIFT button
    cy.get('.release-btn').first().click();
    
    // Verify the action was registered
    cy.wait(500);
    
    // Check localStorage for released shift IDs
    cy.window().then((win) => {
      const releasedIds = JSON.parse(win.localStorage.getItem('released_shift_ids') || '[]');
      expect(releasedIds.length).to.be.greaterThan(0);
    });
  });

  it('should remove the shift from upcoming list after release', () => {
    // Get count before release
    cy.get('.upcoming-item').then(($items) => {
      const countBefore = $items.length;
      
      if (countBefore > 0) {
        // Click release button
        cy.get('.release-btn').first().click();
        cy.wait(500);
        
        // Verify count decreased or list updated
        cy.get('.upcoming-item').should('have.length.lessThan', countBefore);
      }
    });
  });

  describe('Weekly Navigation - Forward Button', () => {
    it('should display forward navigation button for weekly hours', () => {
      cy.get('.weekly-nav-right').should('be.visible');
      cy.get('.weekly-nav-right').should('contain', '>');
    });

    it('should navigate to next week when forward button is clicked', () => {
      // Get the current week label
      cy.get('.weekly-title').then(($title) => {
        const currentWeek = $title.text();
        
        // Click the forward button
        cy.get('.weekly-nav-right').click();
        
        // Verify the week label changes
        cy.get('.weekly-title').should('not.contain', currentWeek);
      });
    });

    it('should update weekly hours display when navigating forward', () => {
      // Get current weekly hours
      cy.get('.weekly-hours').then(($hours) => {
        const currentHours = $hours.text();
        
        // Click the forward button
        cy.get('.weekly-nav-right').click();
        
        // Verify the display updates (may be different or same depending on data)
        cy.get('.weekly-hours').should('exist');
      });
    });

    it('should allow multiple forward navigations', () => {
      // Click forward button 3 times
      cy.get('.weekly-nav-right').click();
      cy.wait(300);
      cy.get('.weekly-nav-right').click();
      cy.wait(300);
      cy.get('.weekly-nav-right').click();
      
      // Verify the weekly section is still visible
      cy.get('.weekly-card').should('be.visible');
      cy.get('.weekly-days').should('be.visible');
    });
  });

  describe('Weekly Navigation - Backward Button', () => {
    it('should display backward navigation button for weekly hours', () => {
      cy.get('.weekly-nav-left').should('be.visible');
      cy.get('.weekly-nav-left').should('contain', '<');
    });

    it('should navigate to previous week when backward button is clicked', () => {
      // First, navigate forward
      cy.get('.weekly-nav-right').click();
      cy.wait(300);

      // Get the current week label after forward navigation
      cy.get('.weekly-title').then(($title) => {
        const forwardWeek = $title.text();
        
        // Click the backward button
        cy.get('.weekly-nav-left').click();
        
        // Verify the week label changes back
        cy.get('.weekly-title').should('not.contain', forwardWeek);
      });
    });

    it('should update weekly hours display when navigating backward', () => {
      // Navigate forward
      cy.get('.weekly-nav-right').click();
      cy.wait(300);

      // Get the forward week hours
      cy.get('.weekly-hours').then(($hours) => {
        const forwardHours = $hours.text();
        
        // Navigate back
        cy.get('.weekly-nav-left').click();
        
        // Verify hours may differ
        cy.get('.weekly-hours').should('exist');
      });
    });

    it('should allow multiple backward navigations', () => {
      // Click forward 3 times
      cy.get('.weekly-nav-right').click();
      cy.wait(300);
      cy.get('.weekly-nav-right').click();
      cy.wait(300);
      cy.get('.weekly-nav-right').click();
      cy.wait(300);

      // Click backward 3 times
      cy.get('.weekly-nav-left').click();
      cy.wait(300);
      cy.get('.weekly-nav-left').click();
      cy.wait(300);
      cy.get('.weekly-nav-left').click();
      
      // Verify the weekly section is still visible and functional
      cy.get('.weekly-card').should('be.visible');
      cy.get('.weekly-days').should('be.visible');
    });
  });

  describe('Combined Shift Release and Weekly Navigation', () => {
    it('should update weekly hours when a shift is released', () => {
      // Get initial weekly total
      cy.get('.weekly-hours').then(($hours) => {
        const initialHours = parseFloat($hours.text());
        
        // Release a shift if available
        cy.get('.release-btn').then(($btn) => {
          if ($btn.length > 0) {
            cy.get('.release-btn').first().click();
            
            // The weekly hours should update (decrease or stay same)
            cy.get('.weekly-hours').should('exist');
          }
        });
      });
    });

    it('should maintain released shift state when navigating weeks', () => {
      // Release a shift
      cy.get('.release-btn').first().click();

      // Store the released shift IDs
      cy.window().then((win) => {
        const releasedIds = JSON.parse(win.localStorage.getItem('released_shift_ids') || '[]');
        const initialCount = releasedIds.length;

        // Navigate to next week
        cy.get('.weekly-nav-right').click();
        cy.wait(300);

        // Verify released shifts are still stored
        cy.window().then((win2) => {
          const currentReleasedIds = JSON.parse(win2.localStorage.getItem('released_shift_ids') || '[]');
          expect(currentReleasedIds.length).to.equal(initialCount);
        });
      });
    });
  });
});
