describe('Supervisor Authentication Flow', () => {
  it('should successfully log in a supervisor and load the dashboard', () => {
    // 1. Start at the landing page
    cy.visit('/home');

    // 2. Select the Supervisor role
    cy.contains('Supervisor').click();

    // 3. Verify routing to the correct login page
    cy.url().should('include', '/login?role=supervisor');
    cy.contains('Supervisor Login').should('be.visible');

    // 4. Input credentials 
    cy.get('input[name="email"]').type('supervisor2@pandaexpress.com');
    cy.get('input[name="password"]').type('pandasupervisor2');

    // 5. Submit the form
    cy.contains('Sign In').click();

    // 6. Verify successful redirect to the dashboard
    cy.url().should('include', '/supervisor-dashboard');

    // 7. Verify the UI hydrated with the user data
    cy.contains('Hello,').should('be.visible');

    // Verify the Schedule View loads the grid
    cy.get('.schedule-grid').should('exist');
    cy.get('.schedule-grid').contains('Schedule view', { matchCase: false }).should('not.exist'); // Ensure it's not totally broken
    
    // Verify Pending Approvals modal opens and buttons render
    cy.contains('Pending approvals').click();
    cy.get('.shift-modal.open').should('exist');
    cy.get('.shift-modal.open').contains('Action Required', { matchCase: false }).should('exist');

    // Verify Overtime / Alerts feature loads
    cy.get('.side-panel').should('exist');
    cy.get('.side-panel').contains('Alerts & tasks', { matchCase: false }).should('exist');
    
    // Wait for the alerts panel to settle after async overtime-risk loading
    cy.get('.side-panel', { timeout: 10000 }).should(($panel) => {
      expect($panel.text()).to.match(/alerts & tasks/i);
      expect($panel.text()).not.to.contain('Checking alerts...');
    });
  });
});