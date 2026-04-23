
import Faq from "./sections/faq";
import AboutUs from "./sections/about";
import Hero from "./sections/hero";
import Services from "./sections/services";
export default function Landing() {

    return (
        <div className="bg-black">
            <Hero />
            <Services />
            <AboutUs />
            <Faq />
        </div>
    );
}
