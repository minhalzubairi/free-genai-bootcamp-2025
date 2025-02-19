## Business Goal:
Architectural diagram(s) is needed that serve as a teaching aid to help stakeholders understand their key components of GenAI workloads. The outcome is to help let stakeholders visualize possible technical paths, technical uncertainty when adopting GenAI for teaching language to students.
We will be guiding key stakeholders through the technical landscape without directly prescribing solutions, while fostering informed discussions about infrastructure choices, integration patterns, and system dependencies across the organization.

## Functional Requirements:
Interest in infrastructure ownership due to the concern of user data privacy and compliances.
Investment of around 15-20K is possible for current active student count of 250 where the students are located in 2 different cities of Japan.

## Non-functional Requirements
Cache hit ratio of at least 70% for common queries
Peak concurrent users should be manageable with auto-scaling groups
All data will be encrypted at rest and in transit
Audit log rotation for at least 30 days
PII data will be detected and masked before reaching the model through guardrails
Single sign-on (SSO) integration likely required for educational institution deployment


## Risks
Infrstructure cost, where the investment allocated would not be enough to spin up the server/infrastructure in-house

## Assumptions:
We would need to make use of a LLM that is Open Source and can be effectively run on hardware that can be set up within the viable investment amount. 
Spin up a server in-house with enough bandwidth to effectively serve current 250 students
Infrastructure costs will scale linearly with usage
System health metrics will be aggregated in a central dashboard
Recovery Time Objective (RTO) of 40 minutes
Activity progress and scores need secure storage

## constraints
Copyrighted materials should not be used at any cost
LLM should be open-source

## data strategy & model selection and development
Knowledge base updates will occur in near real-time
DB will support incremental updates
Data preprocessing pipeline can handle multiple file formats


Regular model evaluation cycles (weekly/bi-weekly)
Input/output guardrails ensure age-appropriate and educational content

Technical Considerations
We can consider any of the following possible open-source models
https://mistral.ai/news/announcing-mistral-7b/
https://huggingface.co/deepseek-ai/deepseek-llm-7b-chat
https://www.microsoft.com/en-us/research/blog/phi-2-the-surprising-power-of-small-language-models/
Base idea is that they should be truly open sourced and preferably run on a consumer grade GPU in-house with good support and documentation.


