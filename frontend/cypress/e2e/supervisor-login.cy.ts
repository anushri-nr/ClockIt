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
  });
});