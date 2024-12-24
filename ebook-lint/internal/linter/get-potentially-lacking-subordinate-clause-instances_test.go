//go:build unit

package linter_test

import (
	"testing"

	"github.com/pjkaufman/go-go-gadgets/ebook-lint/internal/linter"
	"github.com/stretchr/testify/assert"
)

const smallSubordinatingClauseFile = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd">
<html xmlns="http://www.w3.org/1999/xhtml" xml:lang="en" lang="en">
<head>
    <title>Comprehensive Subordinating Clause Test Cases</title>
    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
    <style type="text/css">
        body {
            font-family: Arial, sans-serif;
            line-height: 1.6;
            margin: 2em;
            max-width: 800px;
            background-color: #fafafa;
        }
        h1, h2, h3 {
            color: #333;
            margin-top: 1.5em;
        }
        h1 {
            text-align: center;
            border-bottom: 2px solid #333;
            padding-bottom: 0.5em;
        }
        h2 {
            color: #2c5282;
            border-left: 4px solid #2c5282;
            padding-left: 0.5em;
        }
        h3 {
            color: #4a5568;
            font-style: italic;
        }
        .section {
            margin-bottom: 2em;
            padding: 1.5em;
            border: 1px solid #ddd;
            border-radius: 8px;
            background-color: white;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        .author {
            color: #666;
            font-style: italic;
            margin-bottom: 1em;
            border-bottom: 1px dashed #ccc;
            padding-bottom: 0.5em;
        }
        .date {
            color: #888;
            font-size: 0.9em;
            margin-top: -0.5em;
            margin-bottom: 1em;
        }
        .category {
            background-color: #e2e8f0;
            padding: 0.2em 0.6em;
            border-radius: 4px;
            font-size: 0.9em;
            color: #4a5568;
            display: inline-block;
            margin-bottom: 1em;
        }
        .quote {
            border-left: 3px solid #718096;
            padding-left: 1em;
            margin: 1em 0;
            font-style: italic;
            color: #4a5568;
        }
        .footnote {
            font-size: 0.8em;
            color: #718096;
            margin-top: 1em;
            padding-top: 0.5em;
            border-top: 1px dotted #cbd5e0;
        }
    </style>
</head>
<body>
    <h1>Comprehensive Collection of Academic Essays</h1>
    
    <div class="section">
        <h2>Essay 1: The Impact of Sleep on Academic Performance</h2>
        <p class="author">By Jennifer Martinez</p>
        <p class="date">March 15, 2024</p>
        <p class="category">Academic Research</p>
        
        <p>Sleep plays a crucial role in academic success. Although I was tired, I was not able to sleep. Because of this, but I was unable to do well on the test the next day. The consequences were far-reaching and significant.</p>
        
        <p>Research shows concerning trends in student sleep patterns. While many students prioritize studying, therefore sleep often becomes sacrificed. The balance between academics and rest remains challenging.</p>
        
        <p>Different individuals require varying amounts of sleep. Although eight hours is recommended, however some people function well with less. The quality of sleep matters significantly.</p>
        
        <div class="quote">
            <p>Because sleep affects memory consolidation, thus proper rest becomes essential for learning. Scientific studies continue to support this conclusion.</p>
        </div>
        
        <p class="footnote">Based on research conducted at Sleep Studies Institute, 2024</p>
    </div>

    <div class="section">
        <h2>Essay 2: Environmental Conservation Efforts</h2>
        <p class="author">By Michael Chang</p>
        <p class="date">March 16, 2024</p>
        <p class="category">Environmental Science</p>
        
        <p>Global warming presents unprecedented challenges. While temperatures rise steadily, but immediate action becomes crucial. Scientists worldwide express growing concern.</p>
        
        <p>Conservation efforts require community involvement. Because awareness matters greatly, furthermore education plays a vital role. Local initiatives show promising results.</p>
        
        <p>Renewable energy sources gain importance. Although implementation costs remain high, thus long-term benefits justify investments. Technological advances continue improving efficiency.</p>
        
        <div class="quote">
            <p>While sustainable practices evolve, therefore adaptation becomes necessary. Innovation drives environmental solutions forward.</p>
        </div>
        
        <p class="footnote">Environmental Protection Agency Report Reference, 2024</p>
    </div>

    <div class="section">
        <h2>Essay 3: Digital Technology in Modern Education</h2>
        <p class="author">By Sarah Thompson</p>
        <p class="date">March 17, 2024</p>
        <p class="category">Educational Technology</p>
        
        <p>Technology transforms educational practices. Because digital tools proliferate, however integration challenges persist. Adaptation requires ongoing effort.</p>
        
        <p>Online learning platforms evolve rapidly. While accessibility improves continuously, therefore quality control remains essential. Educational standards must maintain rigor.</p>
        
        <p>Student engagement presents unique challenges. Although digital natives adapt quickly, but proper guidance remains crucial. Teaching methods continue evolving.</p>
        
        <div class="quote">
            <p>Because technology advances rapidly, thus educational methods must adapt. Professional development supports teacher growth.</p>
        </div>
        
        <p class="footnote">Educational Technology Review, Spring 2024</p>
    </div>

    <div class="section">
        <h2>Essay 4: Athletic Performance and Mental Health</h2>
        <p class="author">By James Wilson</p>
        <p class="date">March 18, 2024</p>
        <p class="category">Sports Psychology</p>
        
        <p>Athletic performance depends on mental strength. While physical training matters, furthermore psychological preparation proves essential. Balance maintains peak performance.</p>
        
        <p>Competition pressure affects athletes differently. Although preparation helps significantly, but stress management remains crucial. Support systems play vital roles.</p>
        
        <p>Recovery periods require careful planning. Because rest affects performance, therefore proper scheduling becomes important. Athletes must respect limitations.</p>
        
        <div class="quote">
            <p>While victory motivates greatly, thus maintaining perspective matters most. Mental health deserves priority consideration.</p>
        </div>
        
        <p class="footnote">Sports Psychology Journal, Volume 12, 2024</p>
    </div>

    <div class="section">
        <h2>Essay 5: Modern Literature Analysis</h2>
        <p class="author">By Emma Davidson</p>
        <p class="date">March 19, 2024</p>
        <p class="category">Literary Studies</p>
        
        <p>Contemporary literature reflects societal changes. Because themes evolve constantly, furthermore interpretation requires context. Cultural understanding enhances appreciation.</p>
        
        <p>Digital publishing transforms literature. While traditional methods persist, therefore adaptation becomes necessary. Publishing houses face challenges.</p>
        
        <p>Reader preferences shift continuously. Although classics maintain relevance, but modern formats gain popularity. Digital platforms offer convenience.</p>
        
        <div class="quote">
            <p>Because storytelling evolves, thus narrative structures adapt accordingly. Innovation drives literary development.</p>
        </div>
        
        <p class="footnote">Literary Review Quarterly, Spring 2024</p>
    </div>

    <div class="section">
        <h2>Essay 6: Artificial Intelligence Ethics</h2>
        <p class="author">By Robert Chen</p>
        <p class="date">March 20, 2024</p>
        <p class="category">Technology Ethics</p>
        
        <p>AI development raises ethical questions. While technology advances rapidly, but ethical considerations lag behind. Responsibility requires careful thought.</p>
        
        <p>Privacy concerns grow increasingly important. Because data collection expands, therefore protection becomes crucial. Regulations continue evolving.</p>
        
        <p>Algorithmic bias presents challenges. Although awareness increases steadily, furthermore solutions require development. Research continues advancing.</p>
        
        <div class="quote">
            <p>While innovation drives progress, thus ethical frameworks must adapt. Society bears collective responsibility.</p>
        </div>
        
        <p class="footnote">AI Ethics Journal, Volume 8, 2024</p>
    </div>

    <div class="section">
        <h2>Essay 7: Economic Policy Analysis</h2>
        <p class="author">By Victoria Adams</p>
        <p class="date">March 21, 2024</p>
        <p class="category">Economics</p>
        
        <p>Global economics affects local markets. Because interconnections increase, however independence decreases proportionally. Balance requires careful management.</p>
        
        <p>Market fluctuations challenge stability. While predictions help planning, therefore uncertainty remains constant. Analysis requires ongoing attention.</p>
        
        <p>Policy decisions impact markets significantly. Although research guides decisions, but unexpected factors emerge frequently. Adaptation remains essential.</p>
        
        <div class="quote">
            <p>Because economic factors interrelate, thus comprehensive analysis becomes crucial. Global perspectives enhance understanding.</p>
        </div>
        
        <p class="footnote">Economic Policy Review, March 2024</p>
    </div>

    <div class="section">
        <h2>Essay 8: Healthcare Innovation</h2>
        <p class="author">By Daniel Lee</p>
        <p class="date">March 22, 2024</p>
        <p class="category">Medical Science</p>
        
        <p>Medical technology advances rapidly. While innovation accelerates, furthermore implementation requires caution. Patient safety remains paramount.</p>
        
        <p>Treatment options expand continuously. Because research progresses steadily, therefore possibilities increase exponentially. Ethics guide development.</p>
        
        <p>Healthcare access presents challenges. Although technology improves care, but accessibility requires attention. Equity matters significantly.</p>
        
        <div class="quote">
            <p>While healing advances technically, thus human elements remain essential. Compassion guides healthcare delivery.</p>
        </div>
        
        <p class="footnote">Healthcare Innovation Report, Spring 2024</p>
    </div>

    <div class="section">
        <h2>Essay 9: Urban Development Strategies</h2>
        <p class="author">By Rachel Morgan</p>
        <p class="date">March 23, 2024</p>
        <p class="category">Urban Planning</p>
        
        <p>City planning requires foresight. Because population growth continues, therefore infrastructure needs expansion. Sustainability guides development.</p>
        
        <p>Transportation systems evolve constantly. While efficiency improves steadily, but challenges persist regularly. Innovation drives solutions.</p>
        
        <p>Community needs shape development. Although planning helps significantly, furthermore flexibility remains essential. Adaptation supports growth.</p>
        
        <div class="quote">
            <p>Because cities grow dynamically, thus planning requires vision. Future generations depend upon decisions.</p>
        </div>
        
        <p class="footnote">Urban Development Journal, Volume 15, 2024</p>
    </div>

    <div class="section">
        <h2>Essay 10: Cultural Heritage Preservation</h2>
        <p class="author">By Thomas Anderson</p>
        <p class="date">March 24, 2024</p>
        <p class="category">Cultural Studies</p>
        
        <p>Cultural preservation matters increasingly. While modernization advances rapidly, therefore tradition requires protection. Balance maintains heritage.</p>
        
        <p>Digital archiving transforms preservation. Because technology enables documentation, but authenticity requires verification. Methods continue evolving.</p>
        
        <p>Community involvement supports preservation. Although challenges exist consistently, furthermore participation remains essential. Traditions survive through engagement.</p>
        
        <div class="quote">
            <p>While progress marches forward, thus heritage requires protection. Future generations deserve historical connections.</p>
        </div>
        
        <p class="footnote">Cultural Heritage Review, Spring 2024</p>
    </div>
</body>
</html>`

type getPotentialAlthoughButInstancesTestCase struct {
	inputText           string
	expectedSuggestions map[string]string
}

var getPotentialAlthoughButInstancesTestCases = map[string]getPotentialAlthoughButInstancesTestCase{
	"make sure that a file with no subordinate clause issues does not get a suggestion": {
		inputText: `<p>Here is some content.</p>
<p>Here is some more content</p>`,
		expectedSuggestions: map[string]string{},
	},
	"make sure that a file with a missing comma before an and gets a suggestion": {
		inputText: `<p class="calibre1">This is exactly what Tatsuya was thinking now. Only he was released from school, and Miyuki couldn't miss school with him. </p>
		<p class="calibre1">If it will be a short period, less than a week, this has happened more than once. </p>
		<p class="calibre1"><a id="p57"></a>However, this time there was a chance that this could drag on for a month or more. Although he could "watch" her from afar, but he couldn't stop worrying that Miyuki and Minami would be left alone in this house. </p>
		<p class="calibre1">Maya on the screen looked at Miyuki standing next to Tatsuya, then diagonally behind her from Minami, then looked back at Tatsuya. </p>`,
		expectedSuggestions: map[string]string{
			`		<p class="calibre1"><a id="p57"></a>However, this time there was a chance that this could drag on for a month or more. Although he could "watch" her from afar, but he couldn't stop worrying that Miyuki and Minami would be left alone in this house. </p>`: `		<p class="calibre1"><a id="p57"></a>However, this time there was a chance that this could drag on for a month or more. Although he could "watch" her from afar, he couldn't stop worrying that Miyuki and Minami would be left alone in this house. </p>`,
		},
	},
	"make sure that a paragraph with two potentially broken subordinate clauses with the first not lacking subordination, but the second does should result in a suggestion": {
		inputText: "<p>Sleep plays a crucial role in academic success. Although I was tired, I was not able to sleep. Because of this, but I was unable to do well on the test the next day. The consequences were far-reaching and significant.</p>",
		expectedSuggestions: map[string]string{
			`<p>Sleep plays a crucial role in academic success. Although I was tired, I was not able to sleep. Because of this, but I was unable to do well on the test the next day. The consequences were far-reaching and significant.</p>`: `<p>Sleep plays a crucial role in academic success. Although I was tired, I was not able to sleep. Because of this, I was unable to do well on the test the next day. The consequences were far-reaching and significant.</p>`,
		},
	},
	"Make sure that `Although..., but` results in a suggestion": {
		inputText: "<p>Although the weather was nice, but we stayed inside.</p>",
		expectedSuggestions: map[string]string{
			"<p>Although the weather was nice, but we stayed inside.</p>": "<p>Although the weather was nice, we stayed inside.</p>",
		},
	},
	"Make sure that `Although..., however` results in a suggestion": {
		inputText: "<p>Although she studied hard, however she failed the test.</p>",
		expectedSuggestions: map[string]string{
			"<p>Although she studied hard, however she failed the test.</p>": "<p>Although she studied hard, she failed the test.</p>",
		},
	},
	"Make sure that `Although..., thus` results in a suggestion": {
		inputText: "<p>Although the concert was sold out, thus we watched it online.</p>",
		expectedSuggestions: map[string]string{
			"<p>Although the concert was sold out, thus we watched it online.</p>": "<p>Although the concert was sold out, we watched it online.</p>",
		},
	},
	"Make sure that `Although..., therefore` results in a suggestion": {
		inputText: "<p>Although the evidence was clear, therefore the jury reached a verdict quickly.</p>",
		expectedSuggestions: map[string]string{
			"<p>Although the evidence was clear, therefore the jury reached a verdict quickly.</p>": "<p>Although the evidence was clear, the jury reached a verdict quickly.</p>",
		},
	},
	"Make sure that `Although..., furthermore` results in a suggestion": {
		inputText: "<p>Although the project was challenging, furthermore we completed it on time.</p>",
		expectedSuggestions: map[string]string{
			"<p>Although the project was challenging, furthermore we completed it on time.</p>": "<p>Although the project was challenging, we completed it on time.</p>",
		},
	},
	"Make sure that `While..., therefore` results in a suggestion": {
		inputText: "<p>While the meeting was important, therefore many people attended.</p>",
		expectedSuggestions: map[string]string{
			"<p>While the meeting was important, therefore many people attended.</p>": "<p>While the meeting was important, many people attended.</p>",
		},
	},
	"Make sure that `While..., furthermore` results in a suggestion": {
		inputText: "<p>While the research was ongoing, furthermore new discoveries were made.</p>",
		expectedSuggestions: map[string]string{
			"<p>While the research was ongoing, furthermore new discoveries were made.</p>": "<p>While the research was ongoing, new discoveries were made.</p>",
		},
	},
	"Make sure that `While..., but` results in a suggestion": {
		inputText: "<p>While the car was expensive, but it was worth the investment.</p>",
		expectedSuggestions: map[string]string{
			"<p>While the car was expensive, but it was worth the investment.</p>": "<p>While the car was expensive, it was worth the investment.</p>",
		},
	},
	"Make sure that `While..., however` results in a suggestion": {
		inputText: "<p>While the solution seemed simple, however implementation proved difficult.</p>",
		expectedSuggestions: map[string]string{
			"<p>While the solution seemed simple, however implementation proved difficult.</p>": "<p>While the solution seemed simple, implementation proved difficult.</p>",
		},
	},
	"Make sure that `While..., thus` results in a suggestion": {
		inputText: "<p>While the budget was limited, thus we had to prioritize expenses.</p>",
		expectedSuggestions: map[string]string{
			"<p>While the budget was limited, thus we had to prioritize expenses.</p>": "<p>While the budget was limited, we had to prioritize expenses.</p>",
		},
	},
	"Make sure that `Because..., but` results in a suggestion": {
		inputText: "<p>Because the traffic was heavy, but we arrived late.</p>",
		expectedSuggestions: map[string]string{
			"<p>Because the traffic was heavy, but we arrived late.</p>": "<p>Because the traffic was heavy, we arrived late.</p>",
		},
	},
	"Make sure that `Because..., thus` results in a suggestion": {
		inputText: "<p>Because the deadline was approaching, thus we worked overtime.</p>",
		expectedSuggestions: map[string]string{
			"<p>Because the deadline was approaching, thus we worked overtime.</p>": "<p>Because the deadline was approaching, we worked overtime.</p>",
		},
	},
	"Make sure that `Because..., therefore` results in a suggestion": {
		inputText: "<p>Because the experiment succeeded, therefore we published our findings.</p>",
		expectedSuggestions: map[string]string{
			"<p>Because the experiment succeeded, therefore we published our findings.</p>": "<p>Because the experiment succeeded, we published our findings.</p>",
		},
	},
	"Make sure that `Because..., furthermore` results in a suggestion": {
		inputText: "<p>Because the results were positive, furthermore we continued the research.</p>",
		expectedSuggestions: map[string]string{
			"<p>Because the results were positive, furthermore we continued the research.</p>": "<p>Because the results were positive, we continued the research.</p>",
		},
	},
	"Make sure that `Because..., however` results in a suggestion": {
		inputText: "<p>Because the weather was poor, however the event was cancelled.</p>",
		expectedSuggestions: map[string]string{
			"<p>Because the weather was poor, however the event was cancelled.</p>": "<p>Because the weather was poor, the event was cancelled.</p>",
		},
	},
	"Make sure a file with multiple issues with subordination clauses has them all suggested": {
		inputText: smallSubordinatingClauseFile,
		expectedSuggestions: map[string]string{
			"        <p>Sleep plays a crucial role in academic success. Although I was tired, I was not able to sleep. Because of this, but I was unable to do well on the test the next day. The consequences were far-reaching and significant.</p>": "        <p>Sleep plays a crucial role in academic success. Although I was tired, I was not able to sleep. Because of this, I was unable to do well on the test the next day. The consequences were far-reaching and significant.</p>",
			"        <p>Research shows concerning trends in student sleep patterns. While many students prioritize studying, therefore sleep often becomes sacrificed. The balance between academics and rest remains challenging.</p>":                "        <p>Research shows concerning trends in student sleep patterns. While many students prioritize studying, sleep often becomes sacrificed. The balance between academics and rest remains challenging.</p>",
			"        <p>Different individuals require varying amounts of sleep. Although eight hours is recommended, however some people function well with less. The quality of sleep matters significantly.</p>":                                     "        <p>Different individuals require varying amounts of sleep. Although eight hours is recommended, some people function well with less. The quality of sleep matters significantly.</p>",
			"            <p>Because sleep affects memory consolidation, thus proper rest becomes essential for learning. Scientific studies continue to support this conclusion.</p>":                                                                  "            <p>Because sleep affects memory consolidation, proper rest becomes essential for learning. Scientific studies continue to support this conclusion.</p>",
			"        <p>Global warming presents unprecedented challenges. While temperatures rise steadily, but immediate action becomes crucial. Scientists worldwide express growing concern.</p>":                                                   "        <p>Global warming presents unprecedented challenges. While temperatures rise steadily, immediate action becomes crucial. Scientists worldwide express growing concern.</p>",
			"        <p>Conservation efforts require community involvement. Because awareness matters greatly, furthermore education plays a vital role. Local initiatives show promising results.</p>":                                                "        <p>Conservation efforts require community involvement. Because awareness matters greatly, education plays a vital role. Local initiatives show promising results.</p>",
			"        <p>Renewable energy sources gain importance. Although implementation costs remain high, thus long-term benefits justify investments. Technological advances continue improving efficiency.</p>":                                   "        <p>Renewable energy sources gain importance. Although implementation costs remain high, long-term benefits justify investments. Technological advances continue improving efficiency.</p>",
			"            <p>While sustainable practices evolve, therefore adaptation becomes necessary. Innovation drives environmental solutions forward.</p>":                                                                                        "            <p>While sustainable practices evolve, adaptation becomes necessary. Innovation drives environmental solutions forward.</p>",
			"        <p>Technology transforms educational practices. Because digital tools proliferate, however integration challenges persist. Adaptation requires ongoing effort.</p>":                                                               "        <p>Technology transforms educational practices. Because digital tools proliferate, integration challenges persist. Adaptation requires ongoing effort.</p>",
			"        <p>Online learning platforms evolve rapidly. While accessibility improves continuously, therefore quality control remains essential. Educational standards must maintain rigor.</p>":                                              "        <p>Online learning platforms evolve rapidly. While accessibility improves continuously, quality control remains essential. Educational standards must maintain rigor.</p>",
			"        <p>Student engagement presents unique challenges. Although digital natives adapt quickly, but proper guidance remains crucial. Teaching methods continue evolving.</p>":                                                           "        <p>Student engagement presents unique challenges. Although digital natives adapt quickly, proper guidance remains crucial. Teaching methods continue evolving.</p>",
			"            <p>Because technology advances rapidly, thus educational methods must adapt. Professional development supports teacher growth.</p>":                                                                                           "            <p>Because technology advances rapidly, educational methods must adapt. Professional development supports teacher growth.</p>",
			"        <p>Athletic performance depends on mental strength. While physical training matters, furthermore psychological preparation proves essential. Balance maintains peak performance.</p>":                                             "        <p>Athletic performance depends on mental strength. While physical training matters, psychological preparation proves essential. Balance maintains peak performance.</p>",
			"        <p>Competition pressure affects athletes differently. Although preparation helps significantly, but stress management remains crucial. Support systems play vital roles.</p>":                                                     "        <p>Competition pressure affects athletes differently. Although preparation helps significantly, stress management remains crucial. Support systems play vital roles.</p>",
			"        <p>Recovery periods require careful planning. Because rest affects performance, therefore proper scheduling becomes important. Athletes must respect limitations.</p>":                                                            "        <p>Recovery periods require careful planning. Because rest affects performance, proper scheduling becomes important. Athletes must respect limitations.</p>",
			"            <p>While victory motivates greatly, thus maintaining perspective matters most. Mental health deserves priority consideration.</p>":                                                                                            "            <p>While victory motivates greatly, maintaining perspective matters most. Mental health deserves priority consideration.</p>",
			"        <p>Contemporary literature reflects societal changes. Because themes evolve constantly, furthermore interpretation requires context. Cultural understanding enhances appreciation.</p>":                                           "        <p>Contemporary literature reflects societal changes. Because themes evolve constantly, interpretation requires context. Cultural understanding enhances appreciation.</p>",
			"        <p>Digital publishing transforms literature. While traditional methods persist, therefore adaptation becomes necessary. Publishing houses face challenges.</p>":                                                                   "        <p>Digital publishing transforms literature. While traditional methods persist, adaptation becomes necessary. Publishing houses face challenges.</p>",
			"        <p>Reader preferences shift continuously. Although classics maintain relevance, but modern formats gain popularity. Digital platforms offer convenience.</p>":                                                                     "        <p>Reader preferences shift continuously. Although classics maintain relevance, modern formats gain popularity. Digital platforms offer convenience.</p>",
			"            <p>Because storytelling evolves, thus narrative structures adapt accordingly. Innovation drives literary development.</p>":                                                                                                    "            <p>Because storytelling evolves, narrative structures adapt accordingly. Innovation drives literary development.</p>",
			"        <p>AI development raises ethical questions. While technology advances rapidly, but ethical considerations lag behind. Responsibility requires careful thought.</p>":                                                               "        <p>AI development raises ethical questions. While technology advances rapidly, ethical considerations lag behind. Responsibility requires careful thought.</p>",
			"        <p>Privacy concerns grow increasingly important. Because data collection expands, therefore protection becomes crucial. Regulations continue evolving.</p>":                                                                       "        <p>Privacy concerns grow increasingly important. Because data collection expands, protection becomes crucial. Regulations continue evolving.</p>",
			"        <p>Algorithmic bias presents challenges. Although awareness increases steadily, furthermore solutions require development. Research continues advancing.</p>":                                                                     "        <p>Algorithmic bias presents challenges. Although awareness increases steadily, solutions require development. Research continues advancing.</p>",
			"            <p>While innovation drives progress, thus ethical frameworks must adapt. Society bears collective responsibility.</p>":                                                                                                        "            <p>While innovation drives progress, ethical frameworks must adapt. Society bears collective responsibility.</p>",
			"        <p>Global economics affects local markets. Because interconnections increase, however independence decreases proportionally. Balance requires careful management.</p>":                                                            "        <p>Global economics affects local markets. Because interconnections increase, independence decreases proportionally. Balance requires careful management.</p>",
			"        <p>Market fluctuations challenge stability. While predictions help planning, therefore uncertainty remains constant. Analysis requires ongoing attention.</p>":                                                                    "        <p>Market fluctuations challenge stability. While predictions help planning, uncertainty remains constant. Analysis requires ongoing attention.</p>",
			"        <p>Policy decisions impact markets significantly. Although research guides decisions, but unexpected factors emerge frequently. Adaptation remains essential.</p>":                                                                "        <p>Policy decisions impact markets significantly. Although research guides decisions, unexpected factors emerge frequently. Adaptation remains essential.</p>",
			"            <p>Because economic factors interrelate, thus comprehensive analysis becomes crucial. Global perspectives enhance understanding.</p>":                                                                                         "            <p>Because economic factors interrelate, comprehensive analysis becomes crucial. Global perspectives enhance understanding.</p>",
			"        <p>Medical technology advances rapidly. While innovation accelerates, furthermore implementation requires caution. Patient safety remains paramount.</p>":                                                                         "        <p>Medical technology advances rapidly. While innovation accelerates, implementation requires caution. Patient safety remains paramount.</p>",
			"        <p>Treatment options expand continuously. Because research progresses steadily, therefore possibilities increase exponentially. Ethics guide development.</p>":                                                                    "        <p>Treatment options expand continuously. Because research progresses steadily, possibilities increase exponentially. Ethics guide development.</p>",
			"        <p>Healthcare access presents challenges. Although technology improves care, but accessibility requires attention. Equity matters significantly.</p>":                                                                             "        <p>Healthcare access presents challenges. Although technology improves care, accessibility requires attention. Equity matters significantly.</p>",
			"            <p>While healing advances technically, thus human elements remain essential. Compassion guides healthcare delivery.</p>":                                                                                                      "            <p>While healing advances technically, human elements remain essential. Compassion guides healthcare delivery.</p>",
			"        <p>City planning requires foresight. Because population growth continues, therefore infrastructure needs expansion. Sustainability guides development.</p>":                                                                       "        <p>City planning requires foresight. Because population growth continues, infrastructure needs expansion. Sustainability guides development.</p>",
			"        <p>Transportation systems evolve constantly. While efficiency improves steadily, but challenges persist regularly. Innovation drives solutions.</p>":                                                                              "        <p>Transportation systems evolve constantly. While efficiency improves steadily, challenges persist regularly. Innovation drives solutions.</p>",
			"        <p>Community needs shape development. Although planning helps significantly, furthermore flexibility remains essential. Adaptation supports growth.</p>":                                                                          "        <p>Community needs shape development. Although planning helps significantly, flexibility remains essential. Adaptation supports growth.</p>",
			"            <p>Because cities grow dynamically, thus planning requires vision. Future generations depend upon decisions.</p>":                                                                                                             "            <p>Because cities grow dynamically, planning requires vision. Future generations depend upon decisions.</p>",
			"        <p>Cultural preservation matters increasingly. While modernization advances rapidly, therefore tradition requires protection. Balance maintains heritage.</p>":                                                                    "        <p>Cultural preservation matters increasingly. While modernization advances rapidly, tradition requires protection. Balance maintains heritage.</p>",
			"        <p>Digital archiving transforms preservation. Because technology enables documentation, but authenticity requires verification. Methods continue evolving.</p>":                                                                   "        <p>Digital archiving transforms preservation. Because technology enables documentation, authenticity requires verification. Methods continue evolving.</p>",
			"        <p>Community involvement supports preservation. Although challenges exist consistently, furthermore participation remains essential. Traditions survive through engagement.</p>":                                                  "        <p>Community involvement supports preservation. Although challenges exist consistently, participation remains essential. Traditions survive through engagement.</p>",
			"            <p>While progress marches forward, thus heritage requires protection. Future generations deserve historical connections.</p>":                                                                                                 "            <p>While progress marches forward, heritage requires protection. Future generations deserve historical connections.</p>",
		},
	},
}

func TestGetPotentialAlthoughButInstances(t *testing.T) {
	for name, args := range getPotentialAlthoughButInstancesTestCases {
		t.Run(name, func(t *testing.T) {
			actual := linter.GetPotentiallyLackingSubordinateClauseInstances(args.inputText)

			assert.Equal(t, args.expectedSuggestions, actual)
		})
	}
}
func BenchmarkRegexMatch(b *testing.B) {
	for i := 0; i < b.N; i++ {
		linter.GetPotentiallyLackingSubordinateClauseInstances(smallSubordinatingClauseFile)
	}
}
