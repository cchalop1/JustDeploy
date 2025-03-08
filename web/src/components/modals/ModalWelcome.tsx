import Typewriter from "typewriter-effect";

export default function ModalWelcome() {
  return (
    <div className="fixed w-1/5 bg-white border rounded-lg left-6 font-mono p-4 h-3/4">
      <img src="/hand.png" className="w-28" />
      <div className="font-bold text-2xl">Welcome to your dashoard</div>
      <div className="mt-4">
        <Typewriter
          onInit={(typewriter) => {
            typewriter
              .typeString(
                "Hello, welcome to justdeploy. This is where you'll be able to deploy your project."
              )
              .pauseFor(2500)
              .typeString(
                " Before getting started, you can link a domain name to this server to access your application from the outside."
              )
              .pauseFor(2500)
              .typeString(" Click on settings to specify this domain name.")
              .pauseFor(2500)
              .typeString(
                " Then you can click on the button in the middle of the screen to connect to GitHub and choose a repository to deploy."
              )
              .start();
          }}
          options={{
            delay: 30,
          }}
        />
      </div>
    </div>
  );
}
