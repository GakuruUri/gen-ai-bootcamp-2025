# GenAI Architecting

## Project Overview

This project aims to architect a Generative AI (GenAI) solution for an educational institution with 300 students in Nagasaki. The focus is on privacy, cost control, and scalability.

## Project Requirements

- **Infrastructure**: Own and manage infrastructure to ensure data privacy and control costs.
- **Budget**: Invest in an AI PC with a budget of $10-15K.
- **User Base**: Support 300 active students.

## Project Assumptions

- **Hardware**: Open-source LLMs will run efficiently on the budgeted hardware.
- **Network**: A single server will suffice for bandwidth needs.
- **Model**: Prefer open-source models for transparency and cost control.

## Data Strategy

- **Copyright**: Purchase and store materials to avoid copyright issues.
- **Quality and Security**: Ensure high data quality and implement robust security measures.

## Technical Considerations

### Model Selection

- **Type**: Open-source models like IBM Granite.
- **Deployment**: Self-host for control and cost management.

### Infrastructure

- **Scalability**: Design for future growth.
- **Modularity**: Allow easy updates and replacements.

### Integration

- **APIs**: Develop for seamless integration.
- **CI/CD**: Implement for efficient deployment.

### Monitoring

- **Performance**: Continuous monitoring and feedback loops.
- **KPIs**: Measure business impact.

### Governance

- **Ethical AI**: Policies for responsible AI use.
- **Compliance**: Adhere to regulations and standards.

## Business Considerations

### Use Cases

- **Problem Definition**: Clearly define business problems.
- **Complexity**: Understand integration complexity.
- **Cost**: Identify key cost drivers.
- **Vendor Lock-in**: Avoid with open-source models.

### Guardrails

- **Input/Output**: Manage model inputs and outputs.
- **Caching**: Optimize performance.
- **Sandboxing**: Use containers for testing.