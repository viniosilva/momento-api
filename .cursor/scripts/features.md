a partir da feature-spec mencionada gere os cenários (Gherkin).
agrupe os cenários em comum como o cenário desejavel + cenários de erro.
na pasta desta feature-spec mencionada, para cada grupo (issue) gere um arquivo BDD Feature Scenario (Gherkin + Spec Hybrid) com o formato xx-{feature_name}.md .
analise cada issue criada e identifique o bootstrap necessário para que não haja dependencia entre cada issue. Gere o arquivo 00-initial-setup.md