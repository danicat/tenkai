
import { ValidationRule } from "@/types/domain";

export const ACCEPTANCE_CRITERIA_HEADER = "## Acceptance Criteria";

/**
 * Generates the Markdown list of acceptance criteria based on the provided validation rules.
 */
export function generateAcceptanceCriteria(validation: ValidationRule[]): string {
    if (!validation || validation.length === 0) {
        return "";
    }

    // Filter rules that should be included
    const activeRules = validation.filter(v => v.include_in_prompt);

    if (activeRules.length === 0) {
        return "";
    }

    let criteria = `${ACCEPTANCE_CRITERIA_HEADER}\n`;

    activeRules.forEach(v => {
        if (v.type === 'manual') {
            criteria += `- ${v.target}\n`;
        } else if (v.type === 'test') {
            if (v.min_coverage) {
                criteria += `- Unit tests must pass. You must verify that the *total* project coverage is at least ${v.min_coverage}% by running: \`go test -coverprofile=coverage.out ${v.target || './...'} && go tool cover -func=coverage.out\`\n`;
            } else {
                criteria += `- Unit tests for \`${v.target || './...'}\` must pass.\n`;
            }
        } else if (v.type === 'lint') {
            criteria += `- Linter \`${v.target || './...'}\` must pass`;
            if (v.max_issues !== undefined) criteria += ` with max ${v.max_issues} issues`;
            criteria += '.\n';
        } else if (v.type === 'command') {
            let cmdDesc = `- Command \`${v.command}\``;
            if (v.stdin) cmdDesc += ` with stdin \`${v.stdin}\``;
            cmdDesc += ` must succeed`;

            const requirements: string[] = [];
            if (v.expect_exit_code !== undefined) requirements.push(`exit code ${v.expect_exit_code}`);
            if (v.expect_output) requirements.push(`stdout containing \`${v.expect_output}\``);
            if (v.expect_output_regex) requirements.push(`stdout matching regex \`${v.expect_output_regex}\``);

            if (requirements.length > 0) {
                cmdDesc += ` (must have: ${requirements.join(", ")})`;
            }
            criteria += `${cmdDesc}.\n`;
        } else if (v.type === 'model') {
            criteria += `- AI Review: ${v.prompt || "General code quality check"}\n`;
        }
    });

    return criteria;
}

/**
 * Splits a full prompt into the manual task description and the acceptance criteria section.
 * Useful for loading existing scenarios.
 */
export function splitPrompt(fullPrompt: string): { manual: string, criteria: string } {
    const prompt = fullPrompt || "";
    const headerIndex = prompt.indexOf(ACCEPTANCE_CRITERIA_HEADER);

    if (headerIndex === -1) {
        return { manual: prompt, criteria: "" };
    }

    // We assume the Acceptance Criteria section continues until the next header or end of file.

    // Find where the section ends
    let sectionEndIndex = prompt.length;
    const nextHeaderMatch = prompt.slice(headerIndex + ACCEPTANCE_CRITERIA_HEADER.length).match(/\n##\s/);
    if (nextHeaderMatch && nextHeaderMatch.index !== undefined) {
        sectionEndIndex = headerIndex + ACCEPTANCE_CRITERIA_HEADER.length + nextHeaderMatch.index;
    }

    const manualBefore = prompt.substring(0, headerIndex).trim();
    const manualAfter = prompt.substring(sectionEndIndex).trim();
    const criteriaBlock = prompt.substring(headerIndex, sectionEndIndex);

    // Recombine manual parts
    const manual = manualAfter ? `${manualBefore}\n\n${manualAfter}` : manualBefore;

    return { manual, criteria: criteriaBlock };
}

/**
 * Deprecated: legacy sync function, prefer using split/combine in UI state.
 */
export function syncAcceptanceCriteria(prompt: string, validation: ValidationRule[]): string {
    const { manual } = splitPrompt(prompt);
    const newCriteria = generateAcceptanceCriteria(validation);

    if (!newCriteria) return manual; // No criteria to append

    return manual + "\n\n" + newCriteria;
}
