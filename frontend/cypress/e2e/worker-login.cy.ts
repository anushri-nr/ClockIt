describe('Worker Authentication Flow', () => {
  it('should successfully log in a worker and load the dashboard', () => {
    cy.visit('http://localhost:4200/home');
    cy.contains('Worker').click();
    cy.url().should('include', '/login?role=worker');
    cy.contains('Worker Login').should('be.visible');
    cy.get('input[name="email"]').type('rahul.umamahesha@gmail.com');
    cy.get('input[name="password"]').type('Rahul@2906');
    cy.contains('Sign In').click();
    cy.url().should('include', '/worker-dashboard');
    cy.contains('Hello,').should('be.visible');
  });
});